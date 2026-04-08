package referenceapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"reference-service-go/internal/domain"
	"reference-service-go/internal/outgoing/postgres"
	"reference-service-go/internal/service"
	"strings"

	"github.com/jackc/pgx/v5"
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

var errInt32OutOfRange = errors.New("value is out of int32 range")

// Importer defines the import operations the handler needs.
type Importer interface {
	CreateImport(ctx context.Context, source string) (*domain.Import, error)
	GetImport(ctx context.Context, id string) (*domain.Import, error)
}

// Catcher defines the catch operations the handler needs.
type Catcher interface {
	CreateCatch(ctx context.Context, ballType domain.PokeballType) (*domain.Catch, error)
	GetCatch(ctx context.Context, id string) (*domain.Catch, error)
}

// PokemonReader defines the Pokemon queries the handler needs.
type PokemonReader interface {
	CountPokemon(ctx context.Context) (int64, error)
	CountPokemonByRarity(ctx context.Context, rarity string) (int64, error)
	GetPokemonByID(ctx context.Context, pokedexID int32) (postgres.Pokemon, error)
	ListPokemon(ctx context.Context, params postgres.ListPokemonParams) ([]postgres.Pokemon, error)
	ListPokemonByRarity(ctx context.Context, params postgres.ListPokemonByRarityParams) ([]postgres.Pokemon, error)
}

// APIHandler handles the public service API.
type APIHandler struct {
	logger        *slog.Logger
	importer      Importer
	catcher       Catcher
	pokemonReader PokemonReader
}

// NewHandler creates a new service API handler.
func NewHandler(
	logger *slog.Logger,
	importer Importer,
	catcher Catcher,
	pokemonReader PokemonReader,
) *APIHandler {
	return &APIHandler{
		logger:        logger,
		importer:      importer,
		catcher:       catcher,
		pokemonReader: pokemonReader,
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

	imp, err := h.importer.CreateImport(r.Context(), string(req.Source))
	if err != nil {
		h.logger.ErrorContext(r.Context(), "failed to create import", slog.Any("error", err))
		vital.RespondProblem(r.Context(), w, vital.InternalServerError("failed to create import"))

		return
	}

	resp := ImportResponse{
		Id:        openapi_types.UUID(parseUUIDBytes(imp.ID)),
		Source:    ImportResponseSource(imp.Source),
		Status:    ImportResponseStatus(imp.Status),
		ItemCount: imp.ItemCount,
		CreatedAt: imp.CreatedAt,
		UpdatedAt: imp.UpdatedAt,
	}

	h.logger.InfoContext(r.Context(), "import created",
		slog.String("id", imp.ID),
		slog.String("source", imp.Source),
	)

	w.Header().Set("Location", "/imports/"+imp.ID)
	respondJSON(r.Context(), w, http.StatusCreated, resp, h.logger)
}

// GetImport returns the state of an import by ID.
func (h *APIHandler) GetImport(w http.ResponseWriter, r *http.Request, importID openapi_types.UUID) {
	imp, err := h.importer.GetImport(r.Context(), importID.String())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			vital.RespondProblem(r.Context(), w, vital.NotFound(
				fmt.Sprintf("import %s not found", importID),
			))

			return
		}

		h.logger.ErrorContext(r.Context(), "failed to get import", slog.Any("error", err))
		vital.RespondProblem(r.Context(), w, vital.InternalServerError("failed to get import"))

		return
	}

	respondJSON(r.Context(), w, http.StatusOK, ImportResponse{
		Id:        openapi_types.UUID(parseUUIDBytes(imp.ID)),
		Source:    ImportResponseSource(imp.Source),
		Status:    ImportResponseStatus(imp.Status),
		ItemCount: imp.ItemCount,
		CreatedAt: imp.CreatedAt,
		UpdatedAt: imp.UpdatedAt,
	}, h.logger)
}

// CreateCatch creates and persists a catch.
func (h *APIHandler) CreateCatch(w http.ResponseWriter, r *http.Request) {
	var req CreateCatchRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		vital.RespondProblem(r.Context(), w, vital.BadRequest("invalid request body"))

		return
	}

	catch, err := h.catcher.CreateCatch(r.Context(), domain.PokeballType(req.PokeballType))
	if err != nil {
		if errors.Is(err, service.ErrNoPokemonImported) {
			vital.RespondProblem(r.Context(), w, &vital.ProblemDetail{
				Title:  "No Pokemon Imported",
				Status: http.StatusConflict,
				Detail: "no pokemon have been imported yet, run an import first",
			})

			return
		}

		h.logger.ErrorContext(r.Context(), "failed to create catch", slog.Any("error", err))
		vital.RespondProblem(r.Context(), w, vital.InternalServerError("failed to create catch"))

		return
	}

	resp := CatchResponse{
		Id:           openapi_types.UUID(parseUUIDBytes(catch.ID)),
		Pokemon:      domainToSummary(catch.Pokemon),
		PokeballType: CatchResponsePokeballType(catch.PokeballType),
		IsShiny:      catch.IsShiny,
		CaughtAt:     catch.CaughtAt,
	}

	w.Header().Set("Location", "/catches/"+catch.ID)
	respondJSON(r.Context(), w, http.StatusCreated, resp, h.logger)
}

// GetCatch returns a persisted catch by ID.
func (h *APIHandler) GetCatch(w http.ResponseWriter, r *http.Request, catchID openapi_types.UUID) {
	catch, err := h.catcher.GetCatch(r.Context(), catchID.String())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			vital.RespondProblem(r.Context(), w, vital.NotFound(
				fmt.Sprintf("catch %s not found", catchID),
			))

			return
		}

		h.logger.ErrorContext(r.Context(), "failed to get catch", slog.Any("error", err))
		vital.RespondProblem(r.Context(), w, vital.InternalServerError("failed to get catch"))

		return
	}

	respondJSON(r.Context(), w, http.StatusOK, CatchResponse{
		Id:           openapi_types.UUID(parseUUIDBytes(catch.ID)),
		Pokemon:      domainToSummary(catch.Pokemon),
		PokeballType: CatchResponsePokeballType(catch.PokeballType),
		IsShiny:      catch.IsShiny,
		CaughtAt:     catch.CaughtAt,
	}, h.logger)
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

	items, total, err := h.queryPokemon(r.Context(), params.Rarity, limit, offset)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "failed to list pokemon", slog.Any("error", err))
		vital.RespondProblem(r.Context(), w, vital.InternalServerError("failed to list pokemon"))

		return
	}

	summaries := make([]PokemonSummary, 0, len(items))
	for _, item := range items {
		summaries = append(summaries, dbToSummary(item))
	}

	respondJSON(r.Context(), w, http.StatusOK, PokemonListResponse{
		Items:  summaries,
		Total:  int(total),
		Limit:  limit,
		Offset: offset,
	}, h.logger)
}

// GetPokemon returns a Pokemon by Pokedex ID.
func (h *APIHandler) GetPokemon(w http.ResponseWriter, r *http.Request, pokedexID int) {
	if pokedexID < 0 || pokedexID > maxInt32 {
		vital.RespondProblem(r.Context(), w, vital.BadRequest("pokedex_id is out of range"))

		return
	}

	pokedexID32, err := intToInt32(pokedexID)
	if err != nil {
		vital.RespondProblem(r.Context(), w, vital.BadRequest("pokedex_id is out of range"))

		return
	}

	row, err := h.pokemonReader.GetPokemonByID(r.Context(), pokedexID32)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			vital.RespondProblem(r.Context(), w, vital.NotFound(
				fmt.Sprintf("pokemon %d not found", pokedexID),
			))

			return
		}

		h.logger.ErrorContext(r.Context(), "failed to get pokemon", slog.Any("error", err))
		vital.RespondProblem(r.Context(), w, vital.InternalServerError("failed to get pokemon"))

		return
	}

	respondJSON(r.Context(), w, http.StatusOK, dbToSummary(row), h.logger)
}

func (h *APIHandler) queryPokemon(
	ctx context.Context,
	rarity *ListPokemonParamsRarity,
	limit, offset int,
) ([]postgres.Pokemon, int64, error) {
	if rarity != nil {
		return h.queryPokemonByRarity(ctx, string(*rarity), limit, offset)
	}

	limit32, err := intToInt32(limit)
	if err != nil {
		return nil, 0, fmt.Errorf("converting limit: %w", err)
	}

	offset32, err := intToInt32(offset)
	if err != nil {
		return nil, 0, fmt.Errorf("converting offset: %w", err)
	}

	items, err := h.pokemonReader.ListPokemon(ctx, postgres.ListPokemonParams{
		Limit:  limit32,
		Offset: offset32,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("listing pokemon: %w", err)
	}

	total, err := h.pokemonReader.CountPokemon(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("counting pokemon: %w", err)
	}

	return items, total, nil
}

func (h *APIHandler) queryPokemonByRarity(
	ctx context.Context,
	rarity string,
	limit, offset int,
) ([]postgres.Pokemon, int64, error) {
	limit32, err := intToInt32(limit)
	if err != nil {
		return nil, 0, fmt.Errorf("converting limit: %w", err)
	}

	offset32, err := intToInt32(offset)
	if err != nil {
		return nil, 0, fmt.Errorf("converting offset: %w", err)
	}

	items, err := h.pokemonReader.ListPokemonByRarity(ctx, postgres.ListPokemonByRarityParams{
		Rarity: rarity,
		Limit:  limit32,
		Offset: offset32,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("listing pokemon by rarity: %w", err)
	}

	total, err := h.pokemonReader.CountPokemonByRarity(ctx, rarity)
	if err != nil {
		return nil, 0, fmt.Errorf("counting pokemon by rarity: %w", err)
	}

	return items, total, nil
}

func domainToSummary(pokemon domain.Pokemon) PokemonSummary {
	return PokemonSummary{
		Id:        pokemon.PokedexID,
		Name:      pokemon.Name,
		Rarity:    PokemonSummaryRarity(pokemon.Rarity),
		Types:     pokemon.Types,
		SpriteUrl: pokemon.SpriteURL,
		Stats: PokemonStats{
			Hp:             pokemon.HP,
			Attack:         pokemon.Attack,
			Defense:        pokemon.Defense,
			SpecialAttack:  pokemon.SpecialAttack,
			SpecialDefense: pokemon.SpecialDefense,
			Speed:          pokemon.Speed,
		},
	}
}

func dbToSummary(pokemon postgres.Pokemon) PokemonSummary {
	return PokemonSummary{
		Id:        int(pokemon.PokedexID),
		Name:      pokemon.Name,
		Rarity:    PokemonSummaryRarity(pokemon.Rarity),
		Types:     pokemon.Types,
		SpriteUrl: pokemon.SpriteUrl,
		Stats: PokemonStats{
			Hp:             int(pokemon.Hp),
			Attack:         int(pokemon.Attack),
			Defense:        int(pokemon.Defense),
			SpecialAttack:  int(pokemon.SpecialAttack),
			SpecialDefense: int(pokemon.SpecialDefense),
			Speed:          int(pokemon.Speed),
		},
	}
}

func respondJSON(ctx context.Context, w http.ResponseWriter, status int, body any, logger *slog.Logger) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		logger.ErrorContext(ctx, "failed to encode response", slog.Any("error", err))
	}
}

func parseUUIDBytes(s string) [16]byte {
	var id [16]byte

	var builder strings.Builder

	for _, char := range s {
		if char != '-' {
			builder.WriteRune(char)
		}
	}

	clean := builder.String()

	for index := 0; index < len(id) && index*2+1 < len(clean); index++ {
		_, _ = fmt.Sscanf(clean[index*2:index*2+2], "%02x", &id[index])
	}

	return id
}

func intToInt32(value int) (int32, error) {
	if value < minInt32 || value > maxInt32 {
		return 0, errInt32OutOfRange
	}

	return int32(value), nil
}
