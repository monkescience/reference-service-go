package catch

import (
	"context"
	"errors"
	"reference-service-go/internal/core/pokemon"
	"time"

	"github.com/google/uuid"
)

// ErrNoPokemonImported is returned when no Pokemon have been imported yet.
var (
	ErrNoPokemonImported = errors.New("no pokemon imported yet")
	ErrCatchNotFound     = errors.New("catch not found")
)

// PokeballType represents the type of Pokeball used.
type PokeballType string

const (
	Pokeball   PokeballType = "pokeball"
	GreatBall  PokeballType = "great_ball"
	UltraBall  PokeballType = "ultra_ball"
	MasterBall PokeballType = "master_ball"
)

// Catch represents the result of opening a Pokeball.
type Catch struct {
	ID           uuid.UUID
	Pokemon      pokemon.Pokemon
	PokeballType PokeballType
	IsShiny      bool
	CaughtAt     time.Time
}

// RandomPokemonReader loads a random Pokemon for a rarity.
type RandomPokemonReader interface {
	GetRandomPokemonByRarity(ctx context.Context, rarity pokemon.Rarity) (pokemon.Pokemon, error)
}

// Store persists catches and retrieves them.
type Store interface {
	CreateCatch(ctx context.Context, catch Catch) error
	GetCatch(ctx context.Context, id uuid.UUID) (Catch, error)
}
