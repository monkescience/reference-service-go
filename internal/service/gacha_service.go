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

// OpenPokeball opens a Pokeball and returns a caught Pokemon.
func (s *GachaService) OpenPokeball(ctx context.Context, ballType domain.PokeballType) (*domain.Catch, error) {
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

	s.logger.InfoContext(ctx, "pokeball opened",
		slog.String("pokeball_type", string(ballType)),
		slog.String("pokemon", pokemon.Name),
		slog.String("rarity", string(rarity)),
		slog.Bool("shiny", isShiny),
	)

	return &domain.Catch{
		Pokemon:      pokemon,
		PokeballType: ballType,
		IsShiny:      isShiny,
		CaughtAt:     time.Now(),
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
