package pokemon_test

import (
	"reference-service-go/internal/core/pokemon"
	"testing"

	"github.com/monkescience/testastic"
)

func TestAssignRarity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		isMythical     bool
		isLegendary    bool
		baseExperience int
		want           pokemon.Rarity
	}{
		{
			name:       "mythical overrides other signals",
			isMythical: true,
			want:       pokemon.RarityMythical,
		},
		{
			name:        "legendary beats experience threshold",
			isLegendary: true,
			want:        pokemon.RarityLegendary,
		},
		{
			name:           "rare by base experience",
			baseExperience: pokemon.RareBaseExperienceThreshold,
			want:           pokemon.RarityRare,
		},
		{
			name:           "uncommon by base experience",
			baseExperience: pokemon.UncommonBaseExperienceThreshold,
			want:           pokemon.RarityUncommon,
		},
		{
			name:           "common below uncommon threshold",
			baseExperience: pokemon.UncommonBaseExperienceThreshold - 1,
			want:           pokemon.RarityCommon,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := pokemon.AssignRarity(tt.isMythical, tt.isLegendary, tt.baseExperience)

			testastic.Equal(t, tt.want, got)
		})
	}
}
