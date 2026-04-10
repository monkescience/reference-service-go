package service

import (
	"context"
	"fmt"
	"log/slog"
	"reference-service-go/internal/domain"
	"reference-service-go/internal/outgoing/postgres"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"
)

// PokemonFetcher fetches Pokemon data from an external source.
type PokemonFetcher interface {
	FetchSpeciesCount(ctx context.Context) (int, error)
	FetchPokemon(ctx context.Context, id int) (*domain.Pokemon, error)
}

// ImportService orchestrates Pokemon data imports.
type ImportService struct {
	logger      *slog.Logger
	fetcher     PokemonFetcher
	queries     *postgres.Queries
	pool        *pgxpool.Pool
	concurrency int
	cancelFunc  context.CancelFunc
}

// NewImportService creates a new ImportService.
func NewImportService(
	logger *slog.Logger,
	fetcher PokemonFetcher,
	queries *postgres.Queries,
	pool *pgxpool.Pool,
	concurrency int,
) *ImportService {
	return &ImportService{
		logger:      logger,
		fetcher:     fetcher,
		queries:     queries,
		pool:        pool,
		concurrency: concurrency,
	}
}

// CreateImport creates a new import record and starts async processing.
func (s *ImportService) CreateImport(ctx context.Context, source string) (*domain.Import, error) {
	now := time.Now()

	id, err := newUUIDV7()
	if err != nil {
		return nil, fmt.Errorf("creating import id: %w", err)
	}

	pgID := pgUUIDFromUUID(id)

	imp := domain.Import{
		ID:        id,
		Source:    source,
		Status:    domain.ImportStatusPending,
		ItemCount: 0,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = s.queries.CreateImport(ctx, postgres.CreateImportParams{
		ID:        pgID,
		Source:    source,
		Status:    string(domain.ImportStatusPending),
		ItemCount: 0,
		CreatedAt: pgtype.Timestamptz{Time: now, Valid: true},
		UpdatedAt: pgtype.Timestamptz{Time: now, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("creating import record: %w", err)
	}

	bgCtx, cancel := context.WithCancel(context.Background())
	s.cancelFunc = cancel

	//nolint:contextcheck // Detached from request context so the import outlives the HTTP request.
	go s.runImport(bgCtx, id)

	return &imp, nil
}

// GetImport returns the current state of an import.
func (s *ImportService) GetImport(ctx context.Context, id uuid.UUID) (*domain.Import, error) {
	row, err := s.queries.GetImport(ctx, pgUUIDFromUUID(id))
	if err != nil {
		return nil, fmt.Errorf("getting import: %w", err)
	}

	rowID, err := uuidFromPG(row.ID)
	if err != nil {
		return nil, fmt.Errorf("converting import id: %w", err)
	}

	return &domain.Import{
		ID:        rowID,
		Source:    row.Source,
		Status:    domain.ImportStatus(row.Status),
		ItemCount: int(row.ItemCount),
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}, nil
}

// Shutdown cancels any running imports.
func (s *ImportService) Shutdown() {
	if s.cancelFunc != nil {
		s.cancelFunc()
	}
}

const batchSize = 50

func (s *ImportService) runImport(ctx context.Context, importID uuid.UUID) {
	idStr := importID.String()
	s.logger.Info("starting import", slog.String("import_id", idStr))

	err := s.queries.UpdateImportStatus(ctx, postgres.UpdateImportStatusParams{
		ID:        pgUUIDFromUUID(importID),
		Status:    string(domain.ImportStatusProcessing),
		ItemCount: 0,
	})
	if err != nil {
		s.logger.Error("failed to update import status to processing", slog.Any("error", err))

		return
	}

	count, err := s.fetcher.FetchSpeciesCount(ctx)
	if err != nil {
		s.logger.Error("failed to fetch species count", slog.Any("error", err))
		s.failImport(ctx, importID)

		return
	}

	s.logger.Info("fetching pokemon", slog.Int("total", count))

	pokemon, err := s.fetchAll(ctx, count)
	if err != nil {
		s.logger.Error("failed to fetch pokemon", slog.Any("error", err))
		s.failImport(ctx, importID)

		return
	}

	if !s.upsertAllBatches(ctx, importID, pokemon) {
		return
	}

	s.completeImport(ctx, importID, idStr, len(pokemon))
}

func (s *ImportService) upsertAllBatches(
	ctx context.Context,
	importID uuid.UUID,
	pokemon []domain.Pokemon,
) bool {
	for i := 0; i < len(pokemon); i += batchSize {
		end := min(i+batchSize, len(pokemon))
		batch := pokemon[i:end]

		err := s.upsertBatch(ctx, batch)
		if err != nil {
			s.logger.Error("failed to upsert batch", slog.Any("error", err))
			s.failImport(ctx, importID)

			return false
		}

		//nolint:gosec // Batch index bounded by species count (~1025).
		updateErr := s.queries.UpdateImportStatus(ctx, postgres.UpdateImportStatusParams{
			ID:        pgUUIDFromUUID(importID),
			Status:    string(domain.ImportStatusProcessing),
			ItemCount: int32(end),
		})
		if updateErr != nil {
			s.logger.Error("failed to update item count", slog.Any("error", updateErr))
		}
	}

	return true
}

func (s *ImportService) completeImport(ctx context.Context, importID uuid.UUID, idStr string, count int) {
	//nolint:gosec // Pokemon count bounded by species count (~1025).
	err := s.queries.UpdateImportStatus(ctx, postgres.UpdateImportStatusParams{
		ID:        pgUUIDFromUUID(importID),
		Status:    string(domain.ImportStatusCompleted),
		ItemCount: int32(count),
	})
	if err != nil {
		s.logger.Error("failed to update import status to completed", slog.Any("error", err))
	}

	s.logger.Info("import completed", slog.String("import_id", idStr), slog.Int("count", count))
}

func (s *ImportService) fetchAll(ctx context.Context, count int) ([]domain.Pokemon, error) {
	g, gCtx := errgroup.WithContext(ctx)
	g.SetLimit(s.concurrency)

	results := make(chan domain.Pokemon, count)

	for id := 1; id <= count; id++ {
		g.Go(func() error {
			p, err := s.fetcher.FetchPokemon(gCtx, id)
			if err != nil {
				s.logger.Warn("skipping pokemon",
					slog.Int("id", id),
					slog.Any("error", err),
				)

				return nil
			}

			results <- *p

			return nil
		})
	}

	go func() {
		_ = g.Wait()

		close(results)
	}()

	var pokemon []domain.Pokemon
	for p := range results {
		pokemon = append(pokemon, p)
	}

	err := g.Wait()
	if err != nil {
		return nil, fmt.Errorf("fetching pokemon: %w", err)
	}

	return pokemon, nil
}

func (s *ImportService) upsertBatch(ctx context.Context, pokemon []domain.Pokemon) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}

	defer tx.Rollback(ctx) //nolint:errcheck // Rollback is no-op after commit.

	qtx := s.queries.WithTx(tx)

	for _, p := range pokemon {
		err = qtx.UpsertPokemon(ctx, postgres.UpsertPokemonParams{
			PokedexID:      int32(p.PokedexID), //nolint:gosec // Pokedex IDs are small positive ints.
			Name:           p.Name,
			Rarity:         string(p.Rarity),
			Types:          p.Types,
			SpriteUrl:      p.SpriteURL,
			Hp:             int32(p.HP),             //nolint:gosec // Pokemon stats are small positive ints.
			Attack:         int32(p.Attack),         //nolint:gosec // Pokemon stats are small positive ints.
			Defense:        int32(p.Defense),        //nolint:gosec // Pokemon stats are small positive ints.
			SpecialAttack:  int32(p.SpecialAttack),  //nolint:gosec // Pokemon stats are small positive ints.
			SpecialDefense: int32(p.SpecialDefense), //nolint:gosec // Pokemon stats are small positive ints.
			Speed:          int32(p.Speed),          //nolint:gosec // Pokemon stats are small positive ints.
			BaseExperience: int32(p.BaseExperience), //nolint:gosec // Base experience fits in int32.
			CaptureRate:    int32(p.CaptureRate),    //nolint:gosec // Capture rate is 0-255.
			IsLegendary:    p.IsLegendary,
			IsMythical:     p.IsMythical,
		})
		if err != nil {
			return fmt.Errorf("upserting pokemon %d: %w", p.PokedexID, err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}

func (s *ImportService) failImport(ctx context.Context, importID uuid.UUID) {
	err := s.queries.UpdateImportStatus(ctx, postgres.UpdateImportStatusParams{
		ID:     pgUUIDFromUUID(importID),
		Status: string(domain.ImportStatusFailed),
	})
	if err != nil {
		s.logger.Error("failed to update import status to failed", slog.Any("error", err))
	}
}
