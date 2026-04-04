package domain

import "math/rand/v2"

// ShinyRate is the probability of a catch being shiny (1/512).
const ShinyRate = 1.0 / 512.0

// WeightedTier maps a rarity to its cumulative probability threshold.
type WeightedTier struct {
	Rarity    Rarity
	Threshold float64
}

// OddsTable defines the cumulative probability distribution for each Pokeball type.
//
//nolint:gochecknoglobals // Package-level lookup table for gacha probabilities.
var OddsTable = map[PokeballType][]WeightedTier{
	Pokeball: {
		{RarityCommon, 0.60},
		{RarityUncommon, 0.90},
		{RarityRare, 0.98},
		{RarityLegendary, 0.998},
		{RarityMythical, 1.0},
	},
	GreatBall: {
		{RarityCommon, 0.40},
		{RarityUncommon, 0.75},
		{RarityRare, 0.93},
		{RarityLegendary, 0.99},
		{RarityMythical, 1.0},
	},
	UltraBall: {
		{RarityCommon, 0.20},
		{RarityUncommon, 0.55},
		{RarityRare, 0.85},
		{RarityLegendary, 0.97},
		{RarityMythical, 1.0},
	},
	MasterBall: {
		{RarityCommon, 0.0},
		{RarityUncommon, 0.15},
		{RarityRare, 0.50},
		{RarityLegendary, 0.85},
		{RarityMythical, 1.0},
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
func (DefaultRand) Float64() float64 { return rand.Float64() }

// IntN returns a random int in [0, n).
//
//nolint:gosec // Gacha mechanic does not require crypto randomness.
func (DefaultRand) IntN(n int) int { return rand.IntN(n) }

// RollRarity selects a rarity tier based on the Pokeball's probability distribution.
func RollRarity(ballType PokeballType, rng RandSource) Rarity {
	roll := rng.Float64()
	tiers := OddsTable[ballType]

	for _, tier := range tiers {
		if roll < tier.Threshold {
			return tier.Rarity
		}
	}

	return RarityMythical
}

// RollShiny determines whether a catch is shiny.
func RollShiny(rng RandSource) bool {
	return rng.Float64() < ShinyRate
}
