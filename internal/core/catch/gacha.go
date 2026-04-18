package catch

import (
	"math/rand/v2"
	"reference-service-go/internal/core/pokemon"
)

// ShinyRate is the probability of a catch being shiny (1/512).
const ShinyRate = 1.0 / 512.0

// WeightedTier maps a rarity to its cumulative probability threshold.
type WeightedTier struct {
	Rarity    pokemon.Rarity
	Threshold float64
}

// OddsTable defines the cumulative probability distribution for each Pokeball type.
//
//nolint:gochecknoglobals,mnd // Fixed gameplay odds table shared across the package.
var OddsTable = map[PokeballType][]WeightedTier{
	Pokeball: {
		{Rarity: pokemon.RarityCommon, Threshold: 0.60},
		{Rarity: pokemon.RarityUncommon, Threshold: 0.90},
		{Rarity: pokemon.RarityRare, Threshold: 0.98},
		{Rarity: pokemon.RarityLegendary, Threshold: 0.998},
		{Rarity: pokemon.RarityMythical, Threshold: 1.0},
	},
	GreatBall: {
		{Rarity: pokemon.RarityCommon, Threshold: 0.40},
		{Rarity: pokemon.RarityUncommon, Threshold: 0.75},
		{Rarity: pokemon.RarityRare, Threshold: 0.93},
		{Rarity: pokemon.RarityLegendary, Threshold: 0.99},
		{Rarity: pokemon.RarityMythical, Threshold: 1.0},
	},
	UltraBall: {
		{Rarity: pokemon.RarityCommon, Threshold: 0.20},
		{Rarity: pokemon.RarityUncommon, Threshold: 0.55},
		{Rarity: pokemon.RarityRare, Threshold: 0.85},
		{Rarity: pokemon.RarityLegendary, Threshold: 0.97},
		{Rarity: pokemon.RarityMythical, Threshold: 1.0},
	},
	MasterBall: {
		{Rarity: pokemon.RarityCommon, Threshold: 0.0},
		{Rarity: pokemon.RarityUncommon, Threshold: 0.15},
		{Rarity: pokemon.RarityRare, Threshold: 0.50},
		{Rarity: pokemon.RarityLegendary, Threshold: 0.85},
		{Rarity: pokemon.RarityMythical, Threshold: 1.0},
	},
}

// RandSource abstracts randomness for testability.
type RandSource interface {
	Float64() float64
	IntN(n int) int
}

// DefaultRand wraps math/rand/v2 as the default random source.
type DefaultRand struct{}

// Float64 returns a random float64 in [0.0, 1.0).
//
//nolint:gosec // Gacha mechanic does not require crypto randomness.
func (DefaultRand) Float64() float64 {
	return rand.Float64()
}

// IntN returns a random int in [0, n).
//
//nolint:gosec // Gacha mechanic does not require crypto randomness.
func (DefaultRand) IntN(n int) int {
	return rand.IntN(n)
}

// RollRarity selects a rarity tier based on the Pokeball's probability distribution.
func RollRarity(ballType PokeballType, rng RandSource) pokemon.Rarity {
	roll := rng.Float64()
	tiers := OddsTable[ballType]

	for _, tier := range tiers {
		if roll < tier.Threshold {
			return tier.Rarity
		}
	}

	return pokemon.RarityMythical
}

// RollShiny determines whether a catch is shiny.
func RollShiny(rng RandSource) bool {
	return rng.Float64() < ShinyRate
}
