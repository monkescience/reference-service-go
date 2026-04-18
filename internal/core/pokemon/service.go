package pokemon

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

// Service orchestrates Pokemon imports and catalog queries.
type Service struct {
	fetcher     Fetcher
	imports     ImportStore
	catalog     CatalogStore
	concurrency int
	cancelFunc  context.CancelFunc
}

// NewService creates a new Pokemon service.
func NewService(fetcher Fetcher, imports ImportStore, catalog CatalogStore, concurrency int) *Service {
	return &Service{
		fetcher:     fetcher,
		imports:     imports,
		catalog:     catalog,
		concurrency: concurrency,
	}
}

// CreateImport creates a new import record and starts async processing.
func (s *Service) CreateImport(ctx context.Context, source string) (*Import, error) {
	now := time.Now()

	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("creating import id: %w", err)
	}

	imp := Import{
		ID:        id,
		Source:    source,
		Status:    ImportStatusPending,
		ItemCount: 0,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = s.imports.CreateImport(ctx, imp)
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
func (s *Service) GetImport(ctx context.Context, id uuid.UUID) (*Import, error) {
	imp, err := s.imports.GetImport(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting import: %w", err)
	}

	return &imp, nil
}

// GetPokemonByID returns a Pokemon by Pokedex ID.
func (s *Service) GetPokemonByID(ctx context.Context, pokedexID int) (*Pokemon, error) {
	p, err := s.catalog.GetPokemonByID(ctx, pokedexID)
	if err != nil {
		return nil, fmt.Errorf("getting pokemon: %w", err)
	}

	return &p, nil
}

// ListPokemon returns Pokemon and the matching total count.
func (s *Service) ListPokemon(ctx context.Context, params ListParams) ([]Pokemon, int64, error) {
	items, err := s.catalog.ListPokemon(ctx, params)
	if err != nil {
		return nil, 0, fmt.Errorf("listing pokemon: %w", err)
	}

	total, err := s.catalog.CountPokemon(ctx, params.Rarity)
	if err != nil {
		return nil, 0, fmt.Errorf("counting pokemon: %w", err)
	}

	return items, total, nil
}

// Shutdown cancels any running imports.
func (s *Service) Shutdown() {
	if s.cancelFunc != nil {
		s.cancelFunc()
	}
}

const batchSize = 50

func (s *Service) runImport(ctx context.Context, importID uuid.UUID) {
	idStr := importID.String()
	slog.InfoContext(ctx, "starting import", slog.String("import_id", idStr))

	err := s.imports.UpdateImportStatus(ctx, importID, ImportStatusProcessing, 0)
	if err != nil {
		slog.ErrorContext(ctx, "failed to update import status to processing", slog.Any("error", err))

		return
	}

	count, err := s.fetcher.FetchSpeciesCount(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to fetch species count", slog.Any("error", err))
		s.failImport(ctx, importID)

		return
	}

	slog.InfoContext(ctx, "fetching pokemon", slog.Int("total", count))

	pokemon, err := s.fetchAll(ctx, count)
	if err != nil {
		slog.ErrorContext(ctx, "failed to fetch pokemon", slog.Any("error", err))
		s.failImport(ctx, importID)

		return
	}

	if !s.upsertAllBatches(ctx, importID, pokemon) {
		return
	}

	s.completeImport(ctx, importID, idStr, len(pokemon))
}

func (s *Service) upsertAllBatches(ctx context.Context, importID uuid.UUID, pokemon []Pokemon) bool {
	for i := 0; i < len(pokemon); i += batchSize {
		end := min(i+batchSize, len(pokemon))
		batch := pokemon[i:end]

		err := s.catalog.UpsertPokemonBatch(ctx, batch)
		if err != nil {
			slog.ErrorContext(ctx, "failed to upsert batch", slog.Any("error", err))
			s.failImport(ctx, importID)

			return false
		}

		updateErr := s.imports.UpdateImportStatus(ctx, importID, ImportStatusProcessing, end)
		if updateErr != nil {
			slog.ErrorContext(ctx, "failed to update item count", slog.Any("error", updateErr))
		}
	}

	return true
}

func (s *Service) completeImport(ctx context.Context, importID uuid.UUID, idStr string, count int) {
	err := s.imports.UpdateImportStatus(ctx, importID, ImportStatusCompleted, count)
	if err != nil {
		slog.ErrorContext(ctx, "failed to update import status to completed", slog.Any("error", err))
	}

	slog.InfoContext(ctx, "import completed",
		slog.String("import_id", idStr),
		slog.Int("count", count),
	)
}

func (s *Service) fetchAll(ctx context.Context, count int) ([]Pokemon, error) {
	g, gCtx := errgroup.WithContext(ctx)
	g.SetLimit(s.concurrency)

	results := make(chan Pokemon, count)

	for id := 1; id <= count; id++ {
		pokemonID := id

		g.Go(func() error {
			p, err := s.fetcher.FetchPokemon(gCtx, pokemonID)
			if err != nil {
				slog.WarnContext(gCtx, "skipping pokemon",
					slog.Int("id", pokemonID),
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

	var pokemon []Pokemon
	for p := range results {
		pokemon = append(pokemon, p)
	}

	err := g.Wait()
	if err != nil {
		return nil, fmt.Errorf("fetching pokemon: %w", err)
	}

	return pokemon, nil
}

func (s *Service) failImport(ctx context.Context, importID uuid.UUID) {
	err := s.imports.UpdateImportStatus(ctx, importID, ImportStatusFailed, 0)
	if err != nil {
		slog.ErrorContext(ctx, "failed to update import status to failed", slog.Any("error", err))
	}
}
