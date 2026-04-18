package catch

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// Service handles the Pokeball gacha mechanic.
type Service struct {
	pokemonReader RandomPokemonReader
	store         Store
	rng           RandSource
}

// NewService creates a new catch service.
func NewService(pokemonReader RandomPokemonReader, store Store, rng RandSource) *Service {
	return &Service{
		pokemonReader: pokemonReader,
		store:         store,
		rng:           rng,
	}
}

// CreateCatch creates and persists a catch.
func (s *Service) CreateCatch(ctx context.Context, ballType PokeballType) (*Catch, error) {
	rarity := RollRarity(ballType, s.rng)

	p, err := s.pokemonReader.GetRandomPokemonByRarity(ctx, rarity)
	if err != nil {
		if errors.Is(err, ErrNoPokemonImported) {
			return nil, ErrNoPokemonImported
		}

		return nil, fmt.Errorf("getting random pokemon: %w", err)
	}

	isShiny := RollShiny(s.rng)
	now := time.Now()

	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("creating catch id: %w", err)
	}

	caught := Catch{
		ID:           id,
		Pokemon:      p,
		PokeballType: ballType,
		IsShiny:      isShiny,
		CaughtAt:     now,
	}

	err = s.store.CreateCatch(ctx, caught)
	if err != nil {
		return nil, fmt.Errorf("creating catch: %w", err)
	}

	slog.InfoContext(ctx, "pokeball opened",
		slog.String("pokeball_type", string(ballType)),
		slog.String("pokemon", p.Name),
		slog.String("rarity", string(rarity)),
		slog.Bool("shiny", isShiny),
	)

	return &caught, nil
}

// GetCatch returns a persisted catch by ID.
func (s *Service) GetCatch(ctx context.Context, id uuid.UUID) (*Catch, error) {
	caught, err := s.store.GetCatch(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting catch: %w", err)
	}

	return &caught, nil
}
