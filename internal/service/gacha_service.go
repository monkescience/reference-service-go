package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"reference-service-go/internal/domain"
	"reference-service-go/internal/outgoing/postgres"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// ErrNoPokemonImported is returned when no Pokemon have been imported yet.
var ErrNoPokemonImported = errors.New("no pokemon imported yet")

// GachaService handles the Pokeball gacha mechanic.
type GachaService struct {
	logger  *slog.Logger
	queries *postgres.Queries
	rng     domain.RandSource
}

// NewGachaService creates a new GachaService.
func NewGachaService(logger *slog.Logger, queries *postgres.Queries, rng domain.RandSource) *GachaService {
	return &GachaService{
		logger:  logger,
		queries: queries,
		rng:     rng,
	}
}

// OpenPokeball opens a Pokeball and returns a persisted catch.
func (s *GachaService) OpenPokeball(ctx context.Context, ballType domain.PokeballType) (*domain.Catch, error) {
	return s.CreateCatch(ctx, ballType)
}

// CreateCatch creates and persists a catch.
func (s *GachaService) CreateCatch(ctx context.Context, ballType domain.PokeballType) (*domain.Catch, error) {
	rarity := domain.RollRarity(ballType, s.rng)

	row, err := s.queries.GetRandomPokemonByRarity(ctx, string(rarity))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoPokemonImported
		}

		return nil, fmt.Errorf("getting random pokemon: %w", err)
	}

	pokemon := rowToDomainPokemon(row)
	isShiny := domain.RollShiny(s.rng)
	now := time.Now()
	id := pgtype.UUID{Valid: true}

	copy(id.Bytes[:], newUUIDBytes())

	err = s.queries.CreateCatch(ctx, postgres.CreateCatchParams{
		ID:               id,
		PokemonPokedexID: row.PokedexID,
		PokeballType:     string(ballType),
		IsShiny:          isShiny,
		CaughtAt:         pgtype.Timestamptz{Time: now, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("creating catch: %w", err)
	}

	s.logger.InfoContext(ctx, "pokeball opened",
		slog.String("pokeball_type", string(ballType)),
		slog.String("pokemon", pokemon.Name),
		slog.String("rarity", string(rarity)),
		slog.Bool("shiny", isShiny),
	)

	return &domain.Catch{
		ID:           uuidToString(id),
		Pokemon:      pokemon,
		PokeballType: ballType,
		IsShiny:      isShiny,
		CaughtAt:     now,
	}, nil
}

// GetCatch returns a persisted catch by ID.
func (s *GachaService) GetCatch(ctx context.Context, id string) (*domain.Catch, error) {
	pgID, err := parseUUID(id)
	if err != nil {
		return nil, fmt.Errorf("parsing catch ID: %w", err)
	}

	row, err := s.queries.GetCatch(ctx, pgID)
	if err != nil {
		return nil, fmt.Errorf("getting catch: %w", err)
	}

	return &domain.Catch{
		ID:           uuidToString(row.ID),
		Pokemon:      rowToDomainPokemonFromCatch(row),
		PokeballType: domain.PokeballType(row.PokeballType),
		IsShiny:      row.IsShiny,
		CaughtAt:     row.CaughtAt.Time,
	}, nil
}

func rowToDomainPokemon(row postgres.Pokemon) domain.Pokemon {
	return domain.Pokemon{
		PokedexID:      int(row.PokedexID),
		Name:           row.Name,
		Rarity:         domain.Rarity(row.Rarity),
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

func rowToDomainPokemonFromCatch(row postgres.GetCatchRow) domain.Pokemon {
	return domain.Pokemon{
		PokedexID:      int(row.PokedexID),
		Name:           row.Name,
		Rarity:         domain.Rarity(row.Rarity),
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
