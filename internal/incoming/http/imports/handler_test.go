package importsapi_test

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"reference-service-go/internal/domain"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/monkescience/testastic"

	importsapi "reference-service-go/internal/incoming/http/imports"
)

type fakeImporter struct {
	mu      sync.Mutex
	imports map[string]*domain.Import
}

func newFakeImporter() *fakeImporter {
	return &fakeImporter{
		imports: make(map[string]*domain.Import),
	}
}

func (f *fakeImporter) CreateImport(_ context.Context, source string) (*domain.Import, error) {
	now := time.Now()
	id := generateID()

	imp := &domain.Import{
		ID:        id,
		Source:    source,
		Status:    domain.ImportStatusPending,
		ItemCount: 0,
		CreatedAt: now,
		UpdatedAt: now,
	}

	f.mu.Lock()
	f.imports[id] = imp
	f.mu.Unlock()

	return imp, nil
}

func (f *fakeImporter) GetImport(_ context.Context, id string) (*domain.Import, error) {
	f.mu.Lock()
	imp, ok := f.imports[id]
	f.mu.Unlock()

	if !ok {
		return nil, pgx.ErrNoRows
	}

	return imp, nil
}

const (
	uuidSize        = 16
	uuidVersionByte = 6
	uuidVariantByte = 8
	uuidVersion4    = 0x40
	uuidVersion4Msk = 0x0f
	uuidVariant     = 0x80
	uuidVariantMsk  = 0x3f
)

func generateID() string {
	var id [uuidSize]byte

	_, _ = rand.Read(id[:])

	id[uuidVersionByte] = (id[uuidVersionByte] & uuidVersion4Msk) | uuidVersion4
	id[uuidVariantByte] = (id[uuidVariantByte] & uuidVariantMsk) | uuidVariant

	return fmt.Sprintf("%x-%x-%x-%x-%x", id[0:4], id[4:6], id[6:8], id[8:10], id[10:16])
}

func newTestRouter() http.Handler {
	logger := slog.Default()
	handler := importsapi.NewImportHandler(logger, newFakeImporter())
	router := chi.NewRouter()
	importsapi.HandlerFromMux(handler, router)

	return router
}

func newRequest(method, target string, body string) *http.Request {
	var reader *strings.Reader

	if body != "" {
		reader = strings.NewReader(body)
	}

	if reader != nil {
		return httptest.NewRequestWithContext(context.Background(), method, target, reader)
	}

	return httptest.NewRequestWithContext(context.Background(), method, target, http.NoBody)
}

func loadFixture(t *testing.T, path string) string {
	t.Helper()

	data, err := os.ReadFile(path)
	testastic.NoError(t, err)

	return string(data)
}

func TestCreateImport(t *testing.T) {
	t.Run("returns 201 with valid source", func(t *testing.T) {
		// given: a POST /imports request with a valid source field
		router := newTestRouter()
		req := newRequest(http.MethodPost, "/imports", loadFixture(t, "testdata/create_import/valid_source/request.json"))
		rec := httptest.NewRecorder()

		// when: the request is sent
		router.ServeHTTP(rec, req)

		// then: it responds with 201 and the created import resource
		testastic.Equal(t, http.StatusCreated, rec.Code)
		testastic.Equal(t, "application/json", rec.Header().Get("Content-Type"))
		testastic.AssertJSON(t, "testdata/create_import/valid_source/response.json", rec.Body)
	})

	t.Run("returns 400 with empty body", func(t *testing.T) {
		// given: a POST /imports request with an empty body
		router := newTestRouter()
		req := newRequest(http.MethodPost, "/imports", "")
		rec := httptest.NewRecorder()

		// when: the request is sent
		router.ServeHTTP(rec, req)

		// then: it responds with 400 and a problem detail error
		testastic.Equal(t, http.StatusBadRequest, rec.Code)
		testastic.AssertJSON(t, "testdata/create_import/empty_body/response.json", rec.Body)
	})

	t.Run("returns 400 with missing source", func(t *testing.T) {
		// given: a POST /imports request with an empty JSON object (missing required source)
		router := newTestRouter()
		req := newRequest(http.MethodPost, "/imports", loadFixture(t, "testdata/create_import/missing_source/request.json"))
		rec := httptest.NewRecorder()

		// when: the request is sent
		router.ServeHTTP(rec, req)

		// then: it responds with 400 and a problem detail error
		testastic.Equal(t, http.StatusBadRequest, rec.Code)
		testastic.AssertJSON(t, "testdata/create_import/missing_source/response.json", rec.Body)
	})
}

func TestGetImport(t *testing.T) {
	t.Run("returns 200 for existing import", func(t *testing.T) {
		// given: a previously created import via POST /imports
		router := newTestRouter()

		createReq := newRequest(http.MethodPost, "/imports", loadFixture(t, "testdata/get_import/existing/request.json"))
		createRec := httptest.NewRecorder()
		router.ServeHTTP(createRec, createReq)

		var created importsapi.ImportResponse

		err := json.NewDecoder(createRec.Body).Decode(&created)
		testastic.NoError(t, err)

		// when: GET /imports/{id} is called with the created import's ID
		getReq := newRequest(http.MethodGet, "/imports/"+created.Id.String(), "")
		getRec := httptest.NewRecorder()
		router.ServeHTTP(getRec, getReq)

		// then: it responds with 200 and the import resource
		testastic.Equal(t, http.StatusOK, getRec.Code)
		testastic.Equal(t, "application/json", getRec.Header().Get("Content-Type"))
		testastic.AssertJSON(t, "testdata/get_import/existing/response.json", getRec.Body)
	})

	t.Run("returns 404 for non-existent import", func(t *testing.T) {
		// given: a GET /imports/{id} request with a non-existent UUID
		router := newTestRouter()
		req := newRequest(http.MethodGet, "/imports/550e8400-e29b-41d4-a716-446655440000", "")
		rec := httptest.NewRecorder()

		// when: the request is sent
		router.ServeHTTP(rec, req)

		// then: it responds with 404 and a problem detail error
		testastic.Equal(t, http.StatusNotFound, rec.Code)
		testastic.AssertJSON(t, "testdata/get_import/not_found/response.json", rec.Body)
	})
}
