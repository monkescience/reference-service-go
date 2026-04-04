package pokemonapi

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

	"github.com/jackc/pgx/v5"
	"github.com/monkescience/vital"
)

const (
	defaultLimit  = 20
	defaultOffset = 0
	maxLimit      = 100
)

// PokemonHandler handles Pokemon and Pokeball API requests.
type PokemonHandler struct {
	logger       *slog.Logger
	gachaService *service.GachaService
	queries      *postgres.Queries
}

// NewPokemonHandler creates a new PokemonHandler.
func NewPokemonHandler(
	logger *slog.Logger,
	gachaService *service.GachaService,
	queries *postgres.Queries,
) *PokemonHandler {
	return &PokemonHandler{
		logger:       logger,
		gachaService: gachaService,
		queries:      queries,
	}
}

// OpenPokeball handles POST /pokeball/open.
func (h *PokemonHandler) OpenPokeball(w http.ResponseWriter, r *http.Request) {
	var req OpenPokeballRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		vital.RespondProblem(r.Context(), w, vital.BadRequest("invalid request body"))

		return
	}

	catch, err := h.gachaService.OpenPokeball(r.Context(), domain.PokeballType(req.PokeballType))
	if err != nil {
		if errors.Is(err, service.ErrNoPokemonImported) {
			vital.RespondProblem(r.Context(), w, &vital.ProblemDetail{
				Title:  "No Pokemon Imported",
				Status: http.StatusConflict,
				Detail: "no pokemon have been imported yet, run an import first",
			})

			return
		}

		h.logger.ErrorContext(r.Context(), "failed to open pokeball", slog.Any("error", err))
		vital.RespondProblem(r.Context(), w, vital.InternalServerError("failed to open pokeball"))

		return
	}

	resp := OpenPokeballResponse{
		Pokemon:      domainToSummary(catch.Pokemon),
		PokeballType: OpenPokeballResponsePokeballType(catch.PokeballType),
		IsShiny:      catch.IsShiny,
		CaughtAt:     catch.CaughtAt,
	}

	respondJSON(r.Context(), w, http.StatusOK, resp, h.logger)
}

// ListPokemon handles GET /pokemon.
func (h *PokemonHandler) ListPokemon(w http.ResponseWriter, r *http.Request, params ListPokemonParams) {
	limit := defaultLimit
	if params.Limit != nil {
		limit = *params.Limit
	}

	if limit > maxLimit {
		limit = maxLimit
	}

	offset := defaultOffset
	if params.Offset != nil {
		offset = *params.Offset
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

	resp := PokemonListResponse{
		Items:  summaries,
		Total:  int(total),
		Limit:  limit,
		Offset: offset,
	}

	respondJSON(r.Context(), w, http.StatusOK, resp, h.logger)
}

// GetPokemon handles GET /pokemon/{pokedex_id}.
func (h *PokemonHandler) GetPokemon(w http.ResponseWriter, r *http.Request, pokedexID int) {
	//nolint:gosec // Pokedex IDs are small positive ints.
	row, err := h.queries.GetPokemonByID(r.Context(), int32(pokedexID))
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

func (h *PokemonHandler) queryPokemon(
	ctx context.Context,
	rarity *ListPokemonParamsRarity,
	limit, offset int,
) ([]postgres.Pokemon, int64, error) {
	if rarity != nil {
		return h.queryPokemonByRarity(ctx, string(*rarity), limit, offset)
	}

	//nolint:gosec // Bounded by maxLimit (100) and pagination offset.
	items, err := h.queries.ListPokemon(ctx, postgres.ListPokemonParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("listing pokemon: %w", err)
	}

	total, err := h.queries.CountPokemon(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("counting pokemon: %w", err)
	}

	return items, total, nil
}

func (h *PokemonHandler) queryPokemonByRarity(
	ctx context.Context,
	rarity string,
	limit, offset int,
) ([]postgres.Pokemon, int64, error) {
	//nolint:gosec // Bounded by maxLimit (100) and pagination offset.
	items, err := h.queries.ListPokemonByRarity(ctx, postgres.ListPokemonByRarityParams{
		Rarity: rarity,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("listing pokemon by rarity: %w", err)
	}

	total, err := h.queries.CountPokemonByRarity(ctx, rarity)
	if err != nil {
		return nil, 0, fmt.Errorf("counting pokemon by rarity: %w", err)
	}

	return items, total, nil
}

func domainToSummary(p domain.Pokemon) PokemonSummary {
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

func dbToSummary(p postgres.Pokemon) PokemonSummary {
	return PokemonSummary{
		Id:        int(p.PokedexID),
		Name:      p.Name,
		Rarity:    PokemonSummaryRarity(p.Rarity),
		Types:     p.Types,
		SpriteUrl: p.SpriteUrl,
		Stats: PokemonStats{
			Hp:             int(p.Hp),
			Attack:         int(p.Attack),
			Defense:        int(p.Defense),
			SpecialAttack:  int(p.SpecialAttack),
			SpecialDefense: int(p.SpecialDefense),
			Speed:          int(p.Speed),
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
