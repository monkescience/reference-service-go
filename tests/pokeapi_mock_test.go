//go:build integration

package tests_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
)

type pokeAPIMock struct {
	server *httptest.Server

	mu               sync.RWMutex
	speciesCount     int
	pokemonResponses map[string]string
	speciesResponses map[string]string
}

type pokeAPIMockOption func(t *testing.T, mock *pokeAPIMock)

func newPokeAPIMock(t *testing.T, opts ...pokeAPIMockOption) *pokeAPIMock {
	t.Helper()

	mock := &pokeAPIMock{
		pokemonResponses: make(map[string]string),
		speciesResponses: make(map[string]string),
	}

	for _, opt := range opts {
		opt(t, mock)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v2/pokemon-species", func(w http.ResponseWriter, _ *http.Request) {
		mock.mu.RLock()
		count := mock.speciesCount
		mock.mu.RUnlock()

		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintf(w, `{"count": %d, "results": []}`, count)
	})

	mux.HandleFunc("GET /api/v2/pokemon/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		mock.mu.RLock()
		body, ok := mock.pokemonResponses[id]
		mock.mu.RUnlock()

		if !ok {
			w.WriteHeader(http.StatusNotFound)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, body)
	})

	mux.HandleFunc("GET /api/v2/pokemon-species/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		mock.mu.RLock()
		body, ok := mock.speciesResponses[id]
		mock.mu.RUnlock()

		if !ok {
			w.WriteHeader(http.StatusNotFound)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, body)
	})

	mock.server = httptest.NewServer(mux)
	t.Cleanup(mock.server.Close)

	return mock
}

func withSpeciesCount(count int) pokeAPIMockOption {
	return func(_ *testing.T, mock *pokeAPIMock) {
		mock.mu.Lock()
		mock.speciesCount = count
		mock.mu.Unlock()
	}
}

func withPokemonFixture(id string, pokemonFile string, speciesFile string) pokeAPIMockOption {
	return func(t *testing.T, mock *pokeAPIMock) {
		t.Helper()

		pokemonJSON, err := os.ReadFile(pokemonFile)
		if err != nil {
			t.Fatalf("reading pokemon fixture %s: %v", pokemonFile, err)
		}

		speciesJSON, err := os.ReadFile(speciesFile)
		if err != nil {
			t.Fatalf("reading species fixture %s: %v", speciesFile, err)
		}

		mock.mu.Lock()
		mock.pokemonResponses[id] = string(pokemonJSON)
		mock.speciesResponses[id] = string(speciesJSON)
		mock.mu.Unlock()
	}
}
