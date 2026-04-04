package domain

import "time"

// Rarity represents the rarity tier of a Pokemon.
type Rarity string

const (
	RarityCommon    Rarity = "common"
	RarityUncommon  Rarity = "uncommon"
	RarityRare      Rarity = "rare"
	RarityLegendary Rarity = "legendary"
	RarityMythical  Rarity = "mythical"
)

// PokeballType represents the type of Pokeball used.
type PokeballType string

const (
	Pokeball   PokeballType = "pokeball"
	GreatBall  PokeballType = "great_ball"
	UltraBall  PokeballType = "ultra_ball"
	MasterBall PokeballType = "master_ball"
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

// Catch represents the result of opening a Pokeball.
type Catch struct {
	Pokemon      Pokemon
	PokeballType PokeballType
	IsShiny      bool
	CaughtAt     time.Time
}

// Import represents a Pokemon data import job.
type Import struct {
	ID        string
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
