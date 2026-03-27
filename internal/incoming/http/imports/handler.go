package importsapi

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/monkescience/vital"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// ImportHandler handles import API requests.
type ImportHandler struct {
	logger  *slog.Logger
	mu      sync.Mutex
	imports map[string]ImportResponse
}

// NewImportHandler creates a new import handler.
func NewImportHandler(logger *slog.Logger) *ImportHandler {
	return &ImportHandler{
		logger:  logger,
		imports: make(map[string]ImportResponse),
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

	now := time.Now()
	id := newUUID()

	resp := ImportResponse{
		Id:        id,
		Source:    req.Source,
		Status:    Pending,
		ItemCount: 0,
		CreatedAt: now,
		UpdatedAt: now,
	}

	h.mu.Lock()
	h.imports[id.String()] = resp
	h.mu.Unlock()

	h.logger.InfoContext(r.Context(), "import created",
		slog.String("id", id.String()),
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
	h.mu.Lock()
	resp, ok := h.imports[importID.String()]
	h.mu.Unlock()

	if !ok {
		vital.RespondProblem(r.Context(), w, vital.NotFound(
			fmt.Sprintf("import %s not found", importID),
		))

		return
	}

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "failed to encode response", slog.Any("error", err))
	}
}

// UUID v4 bit masks per RFC 4122.
const (
	uuidVersion4Mask   = 0x0f
	uuidVersion4       = 0x40
	uuidVariantMask    = 0x3f
	uuidVariantRFC4122 = 0x80
)

// newUUID generates a new random UUID v4.
func newUUID() openapi_types.UUID {
	var id [16]byte

	_, _ = rand.Read(id[:])

	id[6] = (id[6] & uuidVersion4Mask) | uuidVersion4
	id[8] = (id[8] & uuidVariantMask) | uuidVariantRFC4122

	return id
}
