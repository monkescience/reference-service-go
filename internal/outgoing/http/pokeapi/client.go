package pokeapi

import (
	"context"
	"fmt"
	"net/http"
	"reference-service-go/internal/domain"
	"strconv"
)

// Fetcher wraps the generated PokeAPI client with domain-level methods.
type Fetcher struct {
	client *ClientWithResponses
}

// NewFetcher creates a new PokeAPI fetcher.
func NewFetcher(httpClient *http.Client, baseURL string) (*Fetcher, error) {
	client, err := NewClientWithResponses(
		baseURL,
		WithHTTPClient(httpClient),
	)
	if err != nil {
		return nil, fmt.Errorf("creating pokeapi client: %w", err)
	}

	return &Fetcher{client: client}, nil
}

// FetchSpeciesCount returns the total number of Pokemon species.
func (f *Fetcher) FetchSpeciesCount(ctx context.Context) (int, error) {
	limit := 0

	resp, err := f.client.ListPokemonSpeciesWithResponse(ctx, &ListPokemonSpeciesParams{
		Limit: &limit,
	})
	if err != nil {
		return 0, fmt.Errorf("fetching species count: %w", err)
	}

	if resp.JSON200 == nil {
		return 0, fmt.Errorf("unexpected status %s from species list", resp.Status()) //nolint:err113 // Dynamic HTTP status.
	}

	return resp.JSON200.Count, nil
}

// FetchPokemon fetches a Pokemon by ID and maps it to a domain Pokemon.
func (f *Fetcher) FetchPokemon(ctx context.Context, id int) (*domain.Pokemon, error) {
	idStr := strconv.Itoa(id)

	pokemonResp, err := f.client.GetPokemonWithResponse(ctx, idStr)
	if err != nil {
		return nil, fmt.Errorf("fetching pokemon %d: %w", id, err)
	}

	if pokemonResp.JSON200 == nil {
		//nolint:err113 // Dynamic HTTP status.
		return nil, fmt.Errorf("unexpected status %s for pokemon %d", pokemonResp.Status(), id)
	}

	speciesResp, err := f.client.GetPokemonSpeciesWithResponse(ctx, idStr)
	if err != nil {
		return nil, fmt.Errorf("fetching species %d: %w", id, err)
	}

	if speciesResp.JSON200 == nil {
		//nolint:err113 // Dynamic HTTP status.
		return nil, fmt.Errorf("unexpected status %s for species %d", speciesResp.Status(), id)
	}

	return mapToDomain(pokemonResp.JSON200, speciesResp.JSON200), nil
}

func mapToDomain(pokemon *PokemonDetail, species *PokemonSpeciesDetail) *domain.Pokemon {
	types := make([]string, 0, len(pokemon.Types))
	for _, t := range pokemon.Types {
		types = append(types, t.Type.Name)
	}

	stats := extractStats(pokemon.Stats)

	baseExp := 0
	if pokemon.BaseExperience != nil {
		baseExp = *pokemon.BaseExperience
	}

	captureRate := 0
	if species.CaptureRate != nil {
		captureRate = *species.CaptureRate
	}

	return &domain.Pokemon{
		PokedexID:      pokemon.Id,
		Name:           pokemon.Name,
		Rarity:         domain.AssignRarity(species.IsMythical, species.IsLegendary, baseExp),
		Types:          types,
		SpriteURL:      selectSprite(pokemon.Sprites),
		HP:             stats.hp,
		Attack:         stats.attack,
		Defense:        stats.defense,
		SpecialAttack:  stats.specialAttack,
		SpecialDefense: stats.specialDefense,
		Speed:          stats.speed,
		BaseExperience: baseExp,
		CaptureRate:    captureRate,
		IsLegendary:    species.IsLegendary,
		IsMythical:     species.IsMythical,
	}
}

type pokemonStats struct {
	hp             int
	attack         int
	defense        int
	specialAttack  int
	specialDefense int
	speed          int
}

func extractStats(stats []PokemonStatEntry) pokemonStats {
	var result pokemonStats

	for _, s := range stats {
		switch s.Stat.Name {
		case "hp":
			result.hp = s.BaseStat
		case "attack":
			result.attack = s.BaseStat
		case "defense":
			result.defense = s.BaseStat
		case "special-attack":
			result.specialAttack = s.BaseStat
		case "special-defense":
			result.specialDefense = s.BaseStat
		case "speed":
			result.speed = s.BaseStat
		}
	}

	return result
}

func selectSprite(sprites PokemonSprites) string {
	if sprites.Other != nil && sprites.Other.OfficialArtwork != nil && sprites.Other.OfficialArtwork.FrontDefault != nil {
		return *sprites.Other.OfficialArtwork.FrontDefault
	}

	if sprites.FrontDefault != nil {
		return *sprites.FrontDefault
	}

	return ""
}
