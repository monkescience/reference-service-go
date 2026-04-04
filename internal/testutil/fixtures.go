package testutil

import (
	"context"
	"reference-service-go/internal/domain"
	"reference-service-go/internal/outgoing/postgres"
	"testing"

	"github.com/monkescience/testastic"
)

const spriteBase = "https://raw.githubusercontent.com/PokeAPI/sprites/" +
	"master/sprites/pokemon/other/official-artwork/"

// Bulbasaur returns a common Pokemon (base_experience=64).
func Bulbasaur() domain.Pokemon {
	spriteURL := spriteBase + "1.png"

	return domain.Pokemon{
		PokedexID:      1,
		Name:           "bulbasaur",
		Rarity:         domain.RarityCommon,
		Types:          []string{"grass", "poison"},
		SpriteURL:      spriteURL,
		HP:             45,
		Attack:         49,
		Defense:        49,
		SpecialAttack:  65,
		SpecialDefense: 65,
		Speed:          45,
		BaseExperience: 64,
		CaptureRate:    45,
		IsLegendary:    false,
		IsMythical:     false,
	}
}

// Pikachu returns an uncommon Pokemon (base_experience=112).
func Pikachu() domain.Pokemon {
	spriteURL := spriteBase + "25.png"

	return domain.Pokemon{
		PokedexID:      25,
		Name:           "pikachu",
		Rarity:         domain.RarityUncommon,
		Types:          []string{"electric"},
		SpriteURL:      spriteURL,
		HP:             35,
		Attack:         55,
		Defense:        40,
		SpecialAttack:  50,
		SpecialDefense: 50,
		Speed:          90,
		BaseExperience: 112,
		CaptureRate:    190,
		IsLegendary:    false,
		IsMythical:     false,
	}
}

// Charizard returns a rare Pokemon (base_experience=267).
func Charizard() domain.Pokemon {
	spriteURL := spriteBase + "6.png"

	return domain.Pokemon{
		PokedexID:      6,
		Name:           "charizard",
		Rarity:         domain.RarityRare,
		Types:          []string{"fire", "flying"},
		SpriteURL:      spriteURL,
		HP:             78,
		Attack:         84,
		Defense:        78,
		SpecialAttack:  109,
		SpecialDefense: 85,
		Speed:          100,
		BaseExperience: 267,
		CaptureRate:    45,
		IsLegendary:    false,
		IsMythical:     false,
	}
}

// Mewtwo returns a legendary Pokemon.
func Mewtwo() domain.Pokemon {
	spriteURL := spriteBase + "150.png"

	return domain.Pokemon{
		PokedexID:      150,
		Name:           "mewtwo",
		Rarity:         domain.RarityLegendary,
		Types:          []string{"psychic"},
		SpriteURL:      spriteURL,
		HP:             106,
		Attack:         110,
		Defense:        90,
		SpecialAttack:  154,
		SpecialDefense: 90,
		Speed:          130,
		BaseExperience: 340,
		CaptureRate:    3,
		IsLegendary:    true,
		IsMythical:     false,
	}
}

// Mew returns a mythical Pokemon.
func Mew() domain.Pokemon {
	spriteURL := spriteBase + "151.png"

	return domain.Pokemon{
		PokedexID:      151,
		Name:           "mew",
		Rarity:         domain.RarityMythical,
		Types:          []string{"psychic"},
		SpriteURL:      spriteURL,
		HP:             100,
		Attack:         100,
		Defense:        100,
		SpecialAttack:  100,
		SpecialDefense: 100,
		Speed:          100,
		BaseExperience: 270,
		CaptureRate:    45,
		IsLegendary:    false,
		IsMythical:     true,
	}
}

// SeedPokemon inserts Pokemon into the database for testing.
func SeedPokemon(ctx context.Context, t *testing.T, queries *postgres.Queries, pokemon ...domain.Pokemon) {
	t.Helper()

	for _, p := range pokemon {
		err := queries.UpsertPokemon(ctx, postgres.UpsertPokemonParams{
			PokedexID:      int32(p.PokedexID), //nolint:gosec // Test data with known small values.
			Name:           p.Name,
			Rarity:         string(p.Rarity),
			Types:          p.Types,
			SpriteUrl:      p.SpriteURL,
			Hp:             int32(p.HP),             //nolint:gosec // Test data with known small values.
			Attack:         int32(p.Attack),         //nolint:gosec // Test data with known small values.
			Defense:        int32(p.Defense),        //nolint:gosec // Test data with known small values.
			SpecialAttack:  int32(p.SpecialAttack),  //nolint:gosec // Test data with known small values.
			SpecialDefense: int32(p.SpecialDefense), //nolint:gosec // Test data with known small values.
			Speed:          int32(p.Speed),          //nolint:gosec // Test data with known small values.
			BaseExperience: int32(p.BaseExperience), //nolint:gosec // Test data with known small values.
			CaptureRate:    int32(p.CaptureRate),    //nolint:gosec // Test data with known small values.
			IsLegendary:    p.IsLegendary,
			IsMythical:     p.IsMythical,
		})
		testastic.NoError(t, err)
	}
}
