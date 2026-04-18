package pokemon

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrImportNotFound  = errors.New("import not found")
	ErrPokemonNotFound = errors.New("pokemon not found")
)

// Rarity represents the rarity tier of a Pokemon.
type Rarity string

const (
	RarityCommon    Rarity = "common"
	RarityUncommon  Rarity = "uncommon"
	RarityRare      Rarity = "rare"
	RarityLegendary Rarity = "legendary"
	RarityMythical  Rarity = "mythical"
)

// Pokemon represents a Pokemon species with its stats and metadata.
type Pokemon struct {
	PokedexID      int
	Name           string
	Rarity         Rarity
	Types          []string
	SpriteURL      string
	HP             int
	Attack         int
	Defense        int
	SpecialAttack  int
	SpecialDefense int
	Speed          int
	BaseExperience int
	CaptureRate    int
	IsLegendary    bool
	IsMythical     bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// Import represents a Pokemon data import job.
type Import struct {
	ID        uuid.UUID
	Source    string
	Status    ImportStatus
	ItemCount int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ImportStatus represents the current state of an import.
type ImportStatus string

const (
	ImportStatusPending    ImportStatus = "pending"
	ImportStatusProcessing ImportStatus = "processing"
	ImportStatusCompleted  ImportStatus = "completed"
	ImportStatusFailed     ImportStatus = "failed"
)

// BaseExperience thresholds for rarity assignment.
const (
	RareBaseExperienceThreshold     = 200
	UncommonBaseExperienceThreshold = 100
)

// ListParams holds catalog query options.
type ListParams struct {
	Rarity *Rarity
	Limit  int
	Offset int
}

// Fetcher fetches Pokemon data from an external source.
type Fetcher interface {
	FetchSpeciesCount(ctx context.Context) (int, error)
	FetchPokemon(ctx context.Context, id int) (*Pokemon, error)
}

// ImportStore persists import state.
type ImportStore interface {
	CreateImport(ctx context.Context, imp Import) error
	GetImport(ctx context.Context, id uuid.UUID) (Import, error)
	UpdateImportStatus(ctx context.Context, id uuid.UUID, status ImportStatus, itemCount int) error
}

// CatalogStore persists and queries Pokemon catalog data.
type CatalogStore interface {
	UpsertPokemonBatch(ctx context.Context, pokemon []Pokemon) error
	GetPokemonByID(ctx context.Context, pokedexID int) (Pokemon, error)
	ListPokemon(ctx context.Context, params ListParams) ([]Pokemon, error)
	CountPokemon(ctx context.Context, rarity *Rarity) (int64, error)
}

// AssignRarity determines a Pokemon's rarity tier based on PokeAPI data.
func AssignRarity(isMythical, isLegendary bool, baseExperience int) Rarity {
	if isMythical {
		return RarityMythical
	}

	if isLegendary {
		return RarityLegendary
	}

	if baseExperience >= RareBaseExperienceThreshold {
		return RarityRare
	}

	if baseExperience >= UncommonBaseExperienceThreshold {
		return RarityUncommon
	}

	return RarityCommon
}
