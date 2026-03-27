package importsapi_test

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
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

func TestCreateImport(t *testing.T) {
	t.Parallel()

	t.Run("returns 201 with valid source", func(t *testing.T) {
		t.Parallel()

		// GIVEN
		router := newTestRouter()
		req := newRequest(http.MethodPost, "/imports", `{"source": "pokeapi"}`)
		rec := httptest.NewRecorder()

		// WHEN
		router.ServeHTTP(rec, req)

		// THEN
		testastic.Equal(t, http.StatusCreated, rec.Code)
		testastic.Equal(t, "application/json", rec.Header().Get("Content-Type"))
		testastic.AssertJSON(t, "testdata/create_import_response.expected.json", rec.Body)
	})

	t.Run("returns 400 with empty body", func(t *testing.T) {
		t.Parallel()

		// GIVEN
		router := newTestRouter()
		req := newRequest(http.MethodPost, "/imports", "")
		rec := httptest.NewRecorder()

		// WHEN
		router.ServeHTTP(rec, req)

		// THEN
		testastic.Equal(t, http.StatusBadRequest, rec.Code)
		testastic.AssertJSON(t, "testdata/create_import_bad_request.expected.json", rec.Body)
	})

	t.Run("returns 400 with missing source", func(t *testing.T) {
		t.Parallel()

		// GIVEN
		router := newTestRouter()
		req := newRequest(http.MethodPost, "/imports", `{}`)
		rec := httptest.NewRecorder()

		// WHEN
		router.ServeHTTP(rec, req)

		// THEN
		testastic.Equal(t, http.StatusBadRequest, rec.Code)
		testastic.AssertJSON(t, "testdata/create_import_bad_request.expected.json", rec.Body)
	})
}

func TestGetImport(t *testing.T) {
	t.Parallel()

	t.Run("returns 200 for existing import", func(t *testing.T) {
		t.Parallel()

		// GIVEN
		router := newTestRouter()

		createReq := newRequest(http.MethodPost, "/imports", `{"source": "pokeapi"}`)
		createRec := httptest.NewRecorder()
		router.ServeHTTP(createRec, createReq)

		var created importsapi.ImportResponse

		err := json.NewDecoder(createRec.Body).Decode(&created)
		testastic.NoError(t, err)

		// WHEN
		getReq := newRequest(http.MethodGet, "/imports/"+created.Id.String(), "")
		getRec := httptest.NewRecorder()
		router.ServeHTTP(getRec, getReq)

		// THEN
		testastic.Equal(t, http.StatusOK, getRec.Code)
		testastic.Equal(t, "application/json", getRec.Header().Get("Content-Type"))
		testastic.AssertJSON(t, "testdata/get_import_response.expected.json", getRec.Body)
	})

	t.Run("returns 404 for non-existent import", func(t *testing.T) {
		t.Parallel()

		// GIVEN
		router := newTestRouter()
		req := newRequest(http.MethodGet, "/imports/550e8400-e29b-41d4-a716-446655440000", "")
		rec := httptest.NewRecorder()

		// WHEN
		router.ServeHTTP(rec, req)

		// THEN
		testastic.Equal(t, http.StatusNotFound, rec.Code)
		testastic.AssertJSON(t, "testdata/get_import_not_found.expected.json", rec.Body)
	})
}
