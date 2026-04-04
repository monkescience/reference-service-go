package importsapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"reference-service-go/internal/domain"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/monkescience/vital"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Importer defines the operations the import handler needs.
type Importer interface {
	CreateImport(ctx context.Context, source string) (*domain.Import, error)
	GetImport(ctx context.Context, id string) (*domain.Import, error)
}

// ImportHandler handles import API requests.
type ImportHandler struct {
	logger   *slog.Logger
	importer Importer
}

// NewImportHandler creates a new import handler.
func NewImportHandler(logger *slog.Logger, importer Importer) *ImportHandler {
	return &ImportHandler{
		logger:   logger,
		importer: importer,
	}
}

// CreateImport triggers a new Pokemon data import.
func (h *ImportHandler) CreateImport(w http.ResponseWriter, r *http.Request) {
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

	imp, err := h.importer.CreateImport(r.Context(), req.Source)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "failed to create import", slog.Any("error", err))
		vital.RespondProblem(r.Context(), w, vital.InternalServerError("failed to create import"))

		return
	}

	resp := ImportResponse{
		Id:        openapi_types.UUID(parseUUIDBytes(imp.ID)),
		Source:    imp.Source,
		Status:    ImportResponseStatus(imp.Status),
		ItemCount: imp.ItemCount,
		CreatedAt: imp.CreatedAt,
		UpdatedAt: imp.UpdatedAt,
	}

	h.logger.InfoContext(r.Context(), "import created",
		slog.String("id", imp.ID),
		slog.String("source", req.Source),
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "failed to encode response", slog.Any("error", err))
	}
}

// GetImport returns the status of an import by ID.
func (h *ImportHandler) GetImport(w http.ResponseWriter, r *http.Request, importID openapi_types.UUID) {
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

	resp := ImportResponse{
		Id:        openapi_types.UUID(parseUUIDBytes(imp.ID)),
		Source:    imp.Source,
		Status:    ImportResponseStatus(imp.Status),
		ItemCount: imp.ItemCount,
		CreatedAt: imp.CreatedAt,
		UpdatedAt: imp.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "failed to encode response", slog.Any("error", err))
	}
}

func parseUUIDBytes(s string) [16]byte {
	var id [16]byte

	var b strings.Builder

	for _, c := range s {
		if c != '-' {
			b.WriteRune(c)
		}
	}

	clean := b.String()

	for i := 0; i < len(id) && i*2+1 < len(clean); i++ {
		_, _ = fmt.Sscanf(clean[i*2:i*2+2], "%02x", &id[i])
	}

	return id
}
