package catch_test

import (
	"reference-service-go/internal/core/catch"
	"reference-service-go/internal/core/pokemon"
	"testing"

	"github.com/monkescience/testastic"
)

type stubRand struct {
	floatValue float64
}

func (s stubRand) Float64() float64 {
	return s.floatValue
}

func (stubRand) IntN(int) int {
	return 0
}

func TestRollRarity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		ball catch.PokeballType
		roll float64
		want pokemon.Rarity
	}{
		{
			name: "pokeball returns common in first bucket",
			ball: catch.Pokeball,
			roll: 0.10,
			want: pokemon.RarityCommon,
		},
		{
			name: "pokeball returns rare in rare bucket",
			ball: catch.Pokeball,
			roll: 0.95,
			want: pokemon.RarityRare,
		},
		{
			name: "masterball returns mythical at upper bound",
			ball: catch.MasterBall,
			roll: 0.99,
			want: pokemon.RarityMythical,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := catch.RollRarity(tt.ball, stubRand{floatValue: tt.roll})

			testastic.Equal(t, tt.want, got)
		})
	}
}

func TestRollShiny(t *testing.T) {
	t.Parallel()

	t.Run("roll below threshold is shiny", func(t *testing.T) {
		t.Parallel()

		isShiny := catch.RollShiny(stubRand{floatValue: catch.ShinyRate / 2})

		testastic.Equal(t, true, isShiny)
	})

	t.Run("roll at threshold is not shiny", func(t *testing.T) {
		t.Parallel()

		isShiny := catch.RollShiny(stubRand{floatValue: catch.ShinyRate})

		testastic.Equal(t, false, isShiny)
	})
}
