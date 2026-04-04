package pokemonapi_test

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"reference-service-go/internal/domain"
	"reference-service-go/internal/outgoing/postgres"
	"reference-service-go/internal/service"
	"reference-service-go/internal/testutil"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/monkescience/testastic"

	pokemonapi "reference-service-go/internal/incoming/http/pokemon"
)

var (
	testPool    *pgxpool.Pool     //nolint:gochecknoglobals // Shared test state via TestMain.
	testQueries *postgres.Queries //nolint:gochecknoglobals // Shared test state via TestMain.
)

func TestMain(m *testing.M) {
	flag.Parse()

	if testing.Short() {
		os.Exit(0)
	}

	ctx := context.Background()

	pg, err := testutil.StartPostgres(ctx)
	if err != nil {
		log.Fatalf("starting postgres: %v", err)
	}

	testPool, testQueries, err = testutil.SetupDatabase(ctx, pg.URL)
	if err != nil {
		log.Fatalf("setting up database: %v", err)
	}

	code := m.Run()

	testPool.Close()

	err = pg.Container.Terminate(ctx)
	if err != nil {
		log.Printf("terminating container: %v", err)
	}

	os.Exit(code)
}

type fakeRand struct {
	floats []float64
	ints   []int
	fi     int
	ii     int
}

func (f *fakeRand) Float64() float64 {
	v := f.floats[f.fi%len(f.floats)]
	f.fi++

	return v
}

func (f *fakeRand) IntN(n int) int {
	v := f.ints[f.ii%len(f.ints)]
	f.ii++

	return v % n
}

func newTestRouter(rng domain.RandSource) http.Handler {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelWarn}))
	gachaService := service.NewGachaService(logger, testQueries, rng)
	handler := pokemonapi.NewPokemonHandler(logger, gachaService, testQueries)
	router := chi.NewRouter()
	pokemonapi.HandlerFromMux(handler, router)

	return router
}

func newRequest(method, target string, body string) *http.Request {
	if body != "" {
		return httptest.NewRequestWithContext(
			context.Background(), method, target, strings.NewReader(body),
		)
	}

	return httptest.NewRequestWithContext(
		context.Background(), method, target, http.NoBody,
	)
}

func TestListPokemon(t *testing.T) {
	ctx := context.Background()

	t.Run("returns paginated list of pokemon", func(t *testing.T) {
		testutil.TruncateTables(ctx, t, testPool)
		testutil.SeedPokemon(ctx, t, testQueries, testutil.Bulbasaur(), testutil.Charizard(), testutil.Pikachu())

		// given: 3 pokemon in the database
		router := newTestRouter(domain.DefaultRand{})
		req := newRequest(http.MethodGet, "/pokemon?limit=2", "") //nolint:contextcheck
		rec := httptest.NewRecorder()

		// when: GET /pokemon with limit=2
		router.ServeHTTP(rec, req)

		// then: returns 200 with 2 items and total=3
		testastic.Equal(t, http.StatusOK, rec.Code)
		testastic.AssertJSON(t, "testdata/list_pokemon/default/response.json", rec.Body)
	})

	t.Run("filters by rarity", func(t *testing.T) {
		testutil.TruncateTables(ctx, t, testPool)
		testutil.SeedPokemon(ctx, t, testQueries, testutil.Bulbasaur(), testutil.Pikachu(), testutil.Charizard())

		// given: pokemon of different rarities
		router := newTestRouter(domain.DefaultRand{})
		req := newRequest(http.MethodGet, "/pokemon?rarity=common", "") //nolint:contextcheck
		rec := httptest.NewRecorder()

		// when: GET /pokemon filtered by common rarity
		router.ServeHTTP(rec, req)

		// then: returns only common pokemon
		testastic.Equal(t, http.StatusOK, rec.Code)
		testastic.AssertJSON(t, "testdata/list_pokemon/by_rarity/response.json", rec.Body)
	})

	t.Run("returns empty list when no pokemon", func(t *testing.T) {
		testutil.TruncateTables(ctx, t, testPool)

		// given: no pokemon in the database
		router := newTestRouter(domain.DefaultRand{})
		req := newRequest(http.MethodGet, "/pokemon", "") //nolint:contextcheck
		rec := httptest.NewRecorder()

		// when: GET /pokemon
		router.ServeHTTP(rec, req)

		// then: returns 200 with empty items
		testastic.Equal(t, http.StatusOK, rec.Code)
		testastic.AssertJSON(t, "testdata/list_pokemon/empty/response.json", rec.Body)
	})
}

func TestGetPokemon(t *testing.T) {
	ctx := context.Background()

	t.Run("returns pokemon by pokedex_id", func(t *testing.T) {
		testutil.TruncateTables(ctx, t, testPool)
		testutil.SeedPokemon(ctx, t, testQueries, testutil.Pikachu())

		// given: Pikachu exists in the database
		router := newTestRouter(domain.DefaultRand{})
		req := newRequest(http.MethodGet, "/pokemon/25", "") //nolint:contextcheck
		rec := httptest.NewRecorder()

		// when: GET /pokemon/25
		router.ServeHTTP(rec, req)

		// then: returns 200 with Pikachu's details
		testastic.Equal(t, http.StatusOK, rec.Code)
		testastic.AssertJSON(t, "testdata/get_pokemon/found/response.json", rec.Body)
	})

	t.Run("returns 404 for non-existent pokemon", func(t *testing.T) {
		testutil.TruncateTables(ctx, t, testPool)

		// given: no pokemon in the database
		router := newTestRouter(domain.DefaultRand{})
		req := newRequest(http.MethodGet, "/pokemon/9999", "") //nolint:contextcheck
		rec := httptest.NewRecorder()

		// when: GET /pokemon/9999
		router.ServeHTTP(rec, req)

		// then: returns 404
		testastic.Equal(t, http.StatusNotFound, rec.Code)
		testastic.AssertJSON(t, "testdata/get_pokemon/not_found/response.json", rec.Body)
	})
}

func TestOpenPokeball(t *testing.T) {
	ctx := context.Background()

	t.Run("returns caught pokemon", func(t *testing.T) {
		testutil.TruncateTables(ctx, t, testPool)
		testutil.SeedPokemon(ctx, t, testQueries, testutil.Bulbasaur())

		// given: a common pokemon exists and RNG is deterministic
		rng := &fakeRand{
			floats: []float64{0.3, 0.5}, // 0.3 → common tier (pokeball), 0.5 → not shiny
			ints:   []int{0},
		}
		router := newTestRouter(rng)
		req := newRequest(http.MethodPost, "/pokeball/open", `{"pokeball_type": "pokeball"}`) //nolint:contextcheck
		rec := httptest.NewRecorder()

		// when: POST /pokeball/open with a pokeball
		router.ServeHTTP(rec, req)

		// then: returns 200 with a caught pokemon
		testastic.Equal(t, http.StatusOK, rec.Code)
		testastic.AssertJSON(t, "testdata/open_pokeball/success/response.json", rec.Body)
	})

	t.Run("returns 400 for invalid body", func(t *testing.T) {
		// given: an empty request body
		router := newTestRouter(domain.DefaultRand{})
		req := newRequest(http.MethodPost, "/pokeball/open", "")
		rec := httptest.NewRecorder()

		// when: POST /pokeball/open with no body
		router.ServeHTTP(rec, req)

		// then: returns 400
		testastic.Equal(t, http.StatusBadRequest, rec.Code)
		testastic.AssertJSON(t, "testdata/open_pokeball/invalid_body/response.json", rec.Body)
	})

	t.Run("returns 409 when no pokemon imported", func(t *testing.T) {
		testutil.TruncateTables(ctx, t, testPool)

		// given: no pokemon in the database
		rng := &fakeRand{
			floats: []float64{0.3, 0.5},
			ints:   []int{0},
		}
		router := newTestRouter(rng)
		req := newRequest(http.MethodPost, "/pokeball/open", `{"pokeball_type": "pokeball"}`) //nolint:contextcheck
		rec := httptest.NewRecorder()

		// when: POST /pokeball/open
		router.ServeHTTP(rec, req)

		// then: returns 409 conflict
		testastic.Equal(t, http.StatusConflict, rec.Code)
		testastic.AssertJSON(t, "testdata/open_pokeball/no_pokemon/response.json", rec.Body)
	})
}
