package importsapi_test

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/monkescience/testastic"

	importsapi "reference-service-go/internal/incoming/http/imports"
)

func newTestRouter() http.Handler {
	logger := slog.Default()
	handler := importsapi.NewImportHandler(logger)
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
	t.Parallel()

	t.Run("returns 201 with valid source", func(t *testing.T) {
		t.Parallel()

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
		t.Parallel()

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
		t.Parallel()

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
	t.Parallel()

	t.Run("returns 200 for existing import", func(t *testing.T) {
		t.Parallel()

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
		t.Parallel()

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
