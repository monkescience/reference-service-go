package referencehttp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"reference-service-go/internal/core/catch"
	"reference-service-go/internal/core/pokemon"

	"github.com/google/uuid"
	"github.com/monkescience/vital"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

const (
	defaultLimit  = 20
	defaultOffset = 0
	maxLimit      = 100
	maxInt32      = int(^uint32(0) >> 1)
	minInt32      = -maxInt32 - 1
)

// PokemonService defines the Pokemon operations the handler needs.
type PokemonService interface {
	CreateImport(ctx context.Context, source string) (*pokemon.Import, error)
	GetImport(ctx context.Context, id uuid.UUID) (*pokemon.Import, error)
	GetPokemonByID(ctx context.Context, pokedexID int) (*pokemon.Pokemon, error)
	ListPokemon(ctx context.Context, params pokemon.ListParams) ([]pokemon.Pokemon, int64, error)
}

// CatchService defines the catch operations the handler needs.
type CatchService interface {
	CreateCatch(ctx context.Context, ballType catch.PokeballType) (*catch.Catch, error)
	GetCatch(ctx context.Context, id uuid.UUID) (*catch.Catch, error)
}

var _ ServerInterface = (*APIHandler)(nil)

// APIHandler handles the public service API.
type APIHandler struct {
	pokemonService PokemonService
	catchService   CatchService
}

// NewHandler creates a new service API handler.
func NewHandler(pokemonService PokemonService, catchService CatchService) *APIHandler {
	return &APIHandler{
		pokemonService: pokemonService,
		catchService:   catchService,
	}
}

// CreateImport creates a new import job.
func (h *APIHandler) CreateImport(w http.ResponseWriter, r *http.Request) {
	var req CreateImportRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		vital.RespondProblem(r.Context(), w, vital.BadRequest("invalid request body"))

		return
	}

	if req.Source == "" {
		vital.RespondProblem(r.Context(), w, vital.BadRequest("source is required"))

		return
	}

	if req.Source != CreateImportRequestSourcePokeapi {
		vital.RespondProblem(r.Context(), w, vital.BadRequest(fmt.Sprintf("unsupported source %q", req.Source)))

		return
	}

	imp, err := h.pokemonService.CreateImport(r.Context(), string(req.Source))
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to create import", slog.Any("error", err))
		vital.RespondProblem(r.Context(), w, vital.InternalServerError("failed to create import"))

		return
	}

	resp := ImportResponse{
		Id:        imp.ID,
		Source:    ImportResponseSource(imp.Source),
		Status:    ImportResponseStatus(imp.Status),
		ItemCount: imp.ItemCount,
		CreatedAt: imp.CreatedAt,
		UpdatedAt: imp.UpdatedAt,
	}

	w.Header().Set("Location", "/imports/"+imp.ID.String())
	respondJSON(r.Context(), w, http.StatusCreated, resp)
}

// GetImport returns the state of an import by ID.
func (h *APIHandler) GetImport(w http.ResponseWriter, r *http.Request, importID openapi_types.UUID) {
	imp, err := h.pokemonService.GetImport(r.Context(), importID)
	if err != nil {
		if errors.Is(err, pokemon.ErrImportNotFound) {
			vital.RespondProblem(r.Context(), w, vital.NotFound(
				fmt.Sprintf("import %s not found", importID),
			))

			return
		}

		slog.ErrorContext(r.Context(), "failed to get import", slog.Any("error", err))
		vital.RespondProblem(r.Context(), w, vital.InternalServerError("failed to get import"))

		return
	}

	respondJSON(r.Context(), w, http.StatusOK, ImportResponse{
		Id:        imp.ID,
		Source:    ImportResponseSource(imp.Source),
		Status:    ImportResponseStatus(imp.Status),
		ItemCount: imp.ItemCount,
		CreatedAt: imp.CreatedAt,
		UpdatedAt: imp.UpdatedAt,
	})
}

// CreateCatch creates and persists a catch.
func (h *APIHandler) CreateCatch(w http.ResponseWriter, r *http.Request) {
	var req CreateCatchRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		vital.RespondProblem(r.Context(), w, vital.BadRequest("invalid request body"))

		return
	}

	caught, err := h.catchService.CreateCatch(r.Context(), catch.PokeballType(req.PokeballType))
	if err != nil {
		if errors.Is(err, catch.ErrNoPokemonImported) {
			vital.RespondProblem(r.Context(), w, &vital.ProblemDetail{
				Title:  "No Pokemon Imported",
				Status: http.StatusConflict,
				Detail: "no pokemon have been imported yet, run an import first",
			})

			return
		}

		slog.ErrorContext(r.Context(), "failed to create catch", slog.Any("error", err))
		vital.RespondProblem(r.Context(), w, vital.InternalServerError("failed to create catch"))

		return
	}

	resp := CatchResponse{
		Id:           caught.ID,
		Pokemon:      pokemonToSummary(caught.Pokemon),
		PokeballType: CatchResponsePokeballType(caught.PokeballType),
		IsShiny:      caught.IsShiny,
		CaughtAt:     caught.CaughtAt,
	}

	w.Header().Set("Location", "/catches/"+caught.ID.String())
	respondJSON(r.Context(), w, http.StatusCreated, resp)
}

// GetCatch returns a persisted catch by ID.
func (h *APIHandler) GetCatch(w http.ResponseWriter, r *http.Request, catchID openapi_types.UUID) {
	caught, err := h.catchService.GetCatch(r.Context(), catchID)
	if err != nil {
		if errors.Is(err, catch.ErrCatchNotFound) {
			vital.RespondProblem(r.Context(), w, vital.NotFound(
				fmt.Sprintf("catch %s not found", catchID),
			))

			return
		}

		slog.ErrorContext(r.Context(), "failed to get catch", slog.Any("error", err))
		vital.RespondProblem(r.Context(), w, vital.InternalServerError("failed to get catch"))

		return
	}

	respondJSON(r.Context(), w, http.StatusOK, CatchResponse{
		Id:           caught.ID,
		Pokemon:      pokemonToSummary(caught.Pokemon),
		PokeballType: CatchResponsePokeballType(caught.PokeballType),
		IsShiny:      caught.IsShiny,
		CaughtAt:     caught.CaughtAt,
	})
}

// ListPokemon lists imported Pokemon.
func (h *APIHandler) ListPokemon(w http.ResponseWriter, r *http.Request, params ListPokemonParams) {
	limit := defaultLimit
	if params.Limit != nil {
		limit = *params.Limit
	}

	if limit < 1 {
		limit = defaultLimit
	}

	if limit > maxLimit {
		limit = maxLimit
	}

	offset := defaultOffset
	if params.Offset != nil {
		offset = *params.Offset
	}

	if offset < 0 {
		offset = defaultOffset
	}

	if offset > maxInt32 {
		offset = maxInt32
	}

	listParams := pokemon.ListParams{Limit: limit, Offset: offset}

	if params.Rarity != nil {
		rarity := pokemon.Rarity(*params.Rarity)
		listParams.Rarity = &rarity
	}

	items, total, err := h.pokemonService.ListPokemon(r.Context(), listParams)
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to list pokemon", slog.Any("error", err))
		vital.RespondProblem(r.Context(), w, vital.InternalServerError("failed to list pokemon"))

		return
	}

	summaries := make([]PokemonSummary, 0, len(items))
	for _, item := range items {
		summaries = append(summaries, pokemonToSummary(item))
	}

	respondJSON(r.Context(), w, http.StatusOK, PokemonListResponse{
		Items:  summaries,
		Total:  int(total),
		Limit:  limit,
		Offset: offset,
	})
}

// GetPokemon returns a Pokemon by Pokedex ID.
func (h *APIHandler) GetPokemon(w http.ResponseWriter, r *http.Request, pokedexID int) {
	if pokedexID < 0 || pokedexID > maxInt32 {
		vital.RespondProblem(r.Context(), w, vital.BadRequest("pokedex_id is out of range"))

		return
	}

	pokemonEntity, err := h.pokemonService.GetPokemonByID(r.Context(), pokedexID)
	if err != nil {
		if errors.Is(err, pokemon.ErrPokemonNotFound) {
			vital.RespondProblem(r.Context(), w, vital.NotFound(
				fmt.Sprintf("pokemon %d not found", pokedexID),
			))

			return
		}

		slog.ErrorContext(r.Context(), "failed to get pokemon", slog.Any("error", err))
		vital.RespondProblem(r.Context(), w, vital.InternalServerError("failed to get pokemon"))

		return
	}

	respondJSON(r.Context(), w, http.StatusOK, pokemonToSummary(*pokemonEntity))
}

func pokemonToSummary(p pokemon.Pokemon) PokemonSummary {
	return PokemonSummary{
		Id:        p.PokedexID,
		Name:      p.Name,
		Rarity:    PokemonSummaryRarity(p.Rarity),
		Types:     p.Types,
		SpriteUrl: p.SpriteURL,
		Stats: PokemonStats{
			Hp:             p.HP,
			Attack:         p.Attack,
			Defense:        p.Defense,
			SpecialAttack:  p.SpecialAttack,
			SpecialDefense: p.SpecialDefense,
			Speed:          p.Speed,
		},
	}
}

func respondJSON(ctx context.Context, w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		slog.ErrorContext(ctx, "failed to encode response", slog.Any("error", err))
	}
}
