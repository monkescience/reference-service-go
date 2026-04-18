package referencepg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reference-service-go/internal/core/catch"
	"reference-service-go/internal/core/pokemon"
	"reference-service-go/internal/outgoing/referencepg/migrations"
	"reference-service-go/internal/outgoing/referencepg/sqlcgen"

	"github.com/exaring/otelpgx"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" // Register pgx as database/sql driver for goose.
	"github.com/pressly/goose/v3"
)

var (
	_ pokemon.ImportStore       = (*Store)(nil)
	_ pokemon.CatalogStore      = (*Store)(nil)
	_ catch.RandomPokemonReader = (*Store)(nil)
	_ catch.Store               = (*Store)(nil)
)

// Store is a PostgreSQL-backed adapter for Pokemon and catch operations.
type Store struct {
	pool    *pgxpool.Pool
	queries *sqlcgen.Queries
}

// New creates a new PostgreSQL store connected to the given DSN.
func New(ctx context.Context, dsn string) (*Store, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse pgx config: %w", err)
	}

	cfg.ConnConfig.Tracer = otelpgx.NewTracer()

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()

		return nil, fmt.Errorf("ping database: %w", err)
	}

	return &Store{
		pool:    pool,
		queries: sqlcgen.New(pool),
	}, nil
}

// Close closes the connection pool.
func (s *Store) Close() {
	s.pool.Close()
}

// CreateImport stores a new import job.
func (s *Store) CreateImport(ctx context.Context, imp pokemon.Import) error {
	err := s.queries.CreateImport(ctx, sqlcgen.CreateImportParams{
		ID:        pgUUIDFromUUID(imp.ID),
		Source:    imp.Source,
		Status:    string(imp.Status),
		ItemCount: int32(imp.ItemCount), //nolint:gosec // Import counts are bounded by species count.
		CreatedAt: pgtype.Timestamptz{Time: imp.CreatedAt, Valid: true},
		UpdatedAt: pgtype.Timestamptz{Time: imp.UpdatedAt, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("create import: %w", err)
	}

	return nil
}

// GetImport returns an import by ID.
func (s *Store) GetImport(ctx context.Context, id uuid.UUID) (pokemon.Import, error) {
	row, err := s.queries.GetImport(ctx, pgUUIDFromUUID(id))
	if err != nil {
		return pokemon.Import{}, fmt.Errorf("get import: %w", err)
	}

	return toCoreImport(row)
}

// UpdateImportStatus updates the status and item count for an import.
func (s *Store) UpdateImportStatus(
	ctx context.Context,
	id uuid.UUID,
	status pokemon.ImportStatus,
	itemCount int,
) error {
	err := s.queries.UpdateImportStatus(ctx, sqlcgen.UpdateImportStatusParams{
		ID:        pgUUIDFromUUID(id),
		Status:    string(status),
		ItemCount: int32(itemCount), //nolint:gosec // Import counts are bounded by species count.
	})
	if err != nil {
		return fmt.Errorf("update import status: %w", err)
	}

	return nil
}

// UpsertPokemonBatch inserts or updates a batch of Pokemon in one transaction.
func (s *Store) UpsertPokemonBatch(ctx context.Context, pokemonBatch []pokemon.Pokemon) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer tx.Rollback(ctx) //nolint:errcheck // Rollback is a no-op after commit.

	queries := s.queries.WithTx(tx)

	for _, p := range pokemonBatch {
		err = queries.UpsertPokemon(ctx, sqlcgen.UpsertPokemonParams{
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
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// GetPokemonByID returns a Pokemon by Pokedex ID.
func (s *Store) GetPokemonByID(ctx context.Context, pokedexID int) (pokemon.Pokemon, error) {
	//nolint:gosec // API validates Pokedex IDs before calling the store.
	row, err := s.queries.GetPokemonByID(ctx, int32(pokedexID))
	if err != nil {
		return pokemon.Pokemon{}, fmt.Errorf("get pokemon by id: %w", err)
	}

	return toCorePokemon(row), nil
}

// ListPokemon returns Pokemon using optional rarity filtering.
func (s *Store) ListPokemon(ctx context.Context, params pokemon.ListParams) ([]pokemon.Pokemon, error) {
	if params.Rarity != nil {
		rows, err := s.queries.ListPokemonByRarity(ctx, sqlcgen.ListPokemonByRarityParams{
			Rarity: string(*params.Rarity),
			Limit:  int32(params.Limit),  //nolint:gosec // Pagination is validated at the API layer.
			Offset: int32(params.Offset), //nolint:gosec // Pagination is validated at the API layer.
		})
		if err != nil {
			return nil, fmt.Errorf("list pokemon by rarity: %w", err)
		}

		return toCorePokemonSlice(rows), nil
	}

	rows, err := s.queries.ListPokemon(ctx, sqlcgen.ListPokemonParams{
		Limit:  int32(params.Limit),  //nolint:gosec // Pagination is validated at the API layer.
		Offset: int32(params.Offset), //nolint:gosec // Pagination is validated at the API layer.
	})
	if err != nil {
		return nil, fmt.Errorf("list pokemon: %w", err)
	}

	return toCorePokemonSlice(rows), nil
}

// CountPokemon returns the total count for the given optional rarity filter.
func (s *Store) CountPokemon(ctx context.Context, rarity *pokemon.Rarity) (int64, error) {
	if rarity != nil {
		count, err := s.queries.CountPokemonByRarity(ctx, string(*rarity))
		if err != nil {
			return 0, fmt.Errorf("count pokemon by rarity: %w", err)
		}

		return count, nil
	}

	count, err := s.queries.CountPokemon(ctx)
	if err != nil {
		return 0, fmt.Errorf("count pokemon: %w", err)
	}

	return count, nil
}

// GetRandomPokemonByRarity returns a random Pokemon for the given rarity.
func (s *Store) GetRandomPokemonByRarity(ctx context.Context, rarity pokemon.Rarity) (pokemon.Pokemon, error) {
	row, err := s.queries.GetRandomPokemonByRarity(ctx, string(rarity))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pokemon.Pokemon{}, catch.ErrNoPokemonImported
		}

		return pokemon.Pokemon{}, fmt.Errorf("get random pokemon by rarity: %w", err)
	}

	return toCorePokemon(row), nil
}

// CreateCatch stores a new catch.
func (s *Store) CreateCatch(ctx context.Context, caught catch.Catch) error {
	err := s.queries.CreateCatch(ctx, sqlcgen.CreateCatchParams{
		ID:               pgUUIDFromUUID(caught.ID),
		PokemonPokedexID: int32(caught.Pokemon.PokedexID), //nolint:gosec // Pokedex IDs are small positive ints.
		PokeballType:     string(caught.PokeballType),
		IsShiny:          caught.IsShiny,
		CaughtAt:         pgtype.Timestamptz{Time: caught.CaughtAt, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("create catch: %w", err)
	}

	return nil
}

// GetCatch returns a persisted catch by ID.
func (s *Store) GetCatch(ctx context.Context, id uuid.UUID) (catch.Catch, error) {
	row, err := s.queries.GetCatch(ctx, pgUUIDFromUUID(id))
	if err != nil {
		return catch.Catch{}, fmt.Errorf("get catch: %w", err)
	}

	catchID, err := uuidFromPG(row.ID)
	if err != nil {
		return catch.Catch{}, fmt.Errorf("convert catch id: %w", err)
	}

	return catch.Catch{
		ID:           catchID,
		Pokemon:      toCorePokemonFromCatch(row),
		PokeballType: catch.PokeballType(row.PokeballType),
		IsShiny:      row.IsShiny,
		CaughtAt:     row.CaughtAt.Time,
	}, nil
}

// Migrate runs all pending goose migrations against the given DSN.
func Migrate(ctx context.Context, dsn string) error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("open db for migrations: %w", err)
	}

	defer db.Close() //nolint:errcheck // Best-effort close after migrations.

	goose.SetBaseFS(migrations.FS)

	err = goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("setting dialect: %w", err)
	}

	err = goose.Up(db, ".")
	if err != nil {
		return fmt.Errorf("running up migrations: %w", err)
	}

	return nil
}

var errNullUUID = errors.New("uuid is null")

func pgUUIDFromUUID(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: [16]byte(id), Valid: true}
}

func uuidFromPG(id pgtype.UUID) (uuid.UUID, error) {
	if !id.Valid {
		return uuid.Nil, errNullUUID
	}

	return uuid.UUID(id.Bytes), nil
}

func toCoreImport(row sqlcgen.Import) (pokemon.Import, error) {
	id, err := uuidFromPG(row.ID)
	if err != nil {
		return pokemon.Import{}, fmt.Errorf("convert import id: %w", err)
	}

	return pokemon.Import{
		ID:        id,
		Source:    row.Source,
		Status:    pokemon.ImportStatus(row.Status),
		ItemCount: int(row.ItemCount),
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}, nil
}

func toCorePokemon(row sqlcgen.Pokemon) pokemon.Pokemon {
	return pokemon.Pokemon{
		PokedexID:      int(row.PokedexID),
		Name:           row.Name,
		Rarity:         pokemon.Rarity(row.Rarity),
		Types:          row.Types,
		SpriteURL:      row.SpriteUrl,
		HP:             int(row.Hp),
		Attack:         int(row.Attack),
		Defense:        int(row.Defense),
		SpecialAttack:  int(row.SpecialAttack),
		SpecialDefense: int(row.SpecialDefense),
		Speed:          int(row.Speed),
		BaseExperience: int(row.BaseExperience),
		CaptureRate:    int(row.CaptureRate),
		IsLegendary:    row.IsLegendary,
		IsMythical:     row.IsMythical,
		CreatedAt:      row.CreatedAt.Time,
		UpdatedAt:      row.UpdatedAt.Time,
	}
}

func toCorePokemonFromCatch(row sqlcgen.GetCatchRow) pokemon.Pokemon {
	return pokemon.Pokemon{
		PokedexID:      int(row.PokedexID),
		Name:           row.Name,
		Rarity:         pokemon.Rarity(row.Rarity),
		Types:          row.Types,
		SpriteURL:      row.SpriteUrl,
		HP:             int(row.Hp),
		Attack:         int(row.Attack),
		Defense:        int(row.Defense),
		SpecialAttack:  int(row.SpecialAttack),
		SpecialDefense: int(row.SpecialDefense),
		Speed:          int(row.Speed),
		BaseExperience: int(row.BaseExperience),
		CaptureRate:    int(row.CaptureRate),
		IsLegendary:    row.IsLegendary,
		IsMythical:     row.IsMythical,
		CreatedAt:      row.CreatedAt.Time,
		UpdatedAt:      row.UpdatedAt.Time,
	}
}

func toCorePokemonSlice(rows []sqlcgen.Pokemon) []pokemon.Pokemon {
	result := make([]pokemon.Pokemon, len(rows))
	for i, row := range rows {
		result[i] = pokemon.Pokemon{
			PokedexID:      int(row.PokedexID),
			Name:           row.Name,
			Rarity:         pokemon.Rarity(row.Rarity),
			Types:          row.Types,
			SpriteURL:      row.SpriteUrl,
			HP:             int(row.Hp),
			Attack:         int(row.Attack),
			Defense:        int(row.Defense),
			SpecialAttack:  int(row.SpecialAttack),
			SpecialDefense: int(row.SpecialDefense),
			Speed:          int(row.Speed),
			BaseExperience: int(row.BaseExperience),
			CaptureRate:    int(row.CaptureRate),
			IsLegendary:    row.IsLegendary,
			IsMythical:     row.IsMythical,
			CreatedAt:      row.CreatedAt.Time,
			UpdatedAt:      row.UpdatedAt.Time,
		}
	}

	return result
}
