package testutil

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

// PokeAPIMock is an httptest server that serves canned PokeAPI responses.
type PokeAPIMock struct {
	Server *httptest.Server

	mu               sync.RWMutex
	speciesCount     int
	pokemonResponses map[string]string
	speciesResponses map[string]string
}

// NewPokeAPIMock creates a mock PokeAPI server.
func NewPokeAPIMock(t *testing.T) *PokeAPIMock {
	t.Helper()

	mock := &PokeAPIMock{
		pokemonResponses: make(map[string]string),
		speciesResponses: make(map[string]string),
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

	mock.Server = httptest.NewServer(mux)
	t.Cleanup(mock.Server.Close)

	return mock
}

// SetSpeciesCount sets the count returned by GET /pokemon-species.
func (m *PokeAPIMock) SetSpeciesCount(count int) {
	m.mu.Lock()
	m.speciesCount = count
	m.mu.Unlock()
}

// AddPokemon registers canned responses for a Pokemon by ID.
func (m *PokeAPIMock) AddPokemon(id string, pokemonJSON string, speciesJSON string) {
	m.mu.Lock()
	m.pokemonResponses[id] = pokemonJSON
	m.speciesResponses[id] = speciesJSON
	m.mu.Unlock()
}
