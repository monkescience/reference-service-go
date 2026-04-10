//go:build integration

package tests_test

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/monkescience/testastic"
)

func newScenarioPokeAPIMock(t *testing.T, fixtureDir string) *pokeAPIMock {
	t.Helper()

	return newPokeAPIMock(t,
		withSpeciesCount(2),
		withPokemonFixture("1",
			fixtureDir+"/pokeapi_first_pokemon.json",
			fixtureDir+"/pokeapi_first_species.json",
		),
		withPokemonFixture("2",
			fixtureDir+"/pokeapi_second_pokemon.json",
			fixtureDir+"/pokeapi_second_species.json",
		),
	)
}

func newCatchAfterImportMock(t *testing.T) *pokeAPIMock {
	t.Helper()

	fixtureDir := "testdata/open_pokeball_after_import"

	return newPokeAPIMock(t,
		withSpeciesCount(5),
		withPokemonFixture("1", fixtureDir+"/pokeapi_first_pokemon.json", fixtureDir+"/pokeapi_first_species.json"),
		withPokemonFixture("2", fixtureDir+"/pokeapi_second_pokemon.json", fixtureDir+"/pokeapi_second_species.json"),
		withPokemonFixture("3", fixtureDir+"/pokeapi_third_pokemon.json", fixtureDir+"/pokeapi_third_species.json"),
		withPokemonFixture("4", fixtureDir+"/pokeapi_fourth_pokemon.json", fixtureDir+"/pokeapi_fourth_species.json"),
		withPokemonFixture("5", fixtureDir+"/pokeapi_fifth_pokemon.json", fixtureDir+"/pokeapi_fifth_species.json"),
	)
}

func importPokemonForSetup(t *testing.T, procURL string) {
	t.Helper()

	resp := doPost(t, procURL+"/imports", `{"source": "pokeapi"}`)
	testastic.Equal(t, http.StatusCreated, resp.StatusCode)
	body := readBody(t, resp)

	var importResp createdImportResponse

	decodeJSON(t, body, &importResp)

	testastic.EventuallyEqual(t, "completed", func() string {
		resp := doGet(t, procURL+"/imports/"+importResp.ID)
		body := readBody(t, resp)

		var status importStatusResponse

		decodeJSON(t, body, &status)

		if status.Status == "failed" {
			t.Fatalf("import %s failed", importResp.ID)
		}

		return status.Status
	}, 30*time.Second)
}

func TestHealthEndpoint(t *testing.T) {
	// given: a running service with a fresh PokeAPI fake
	mock := newPokeAPIMock(t)
	proc := startService(t, mock.server.URL+"/api/v2")

	// when: GET /health/live is called
	resp := doGet(t, proc.URL()+"/health/live")

	// then: the live endpoint responds successfully
	testastic.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestImportFlow(t *testing.T) {
	// given: a running service and a PokeAPI fake serving two pokemon
	mock := newScenarioPokeAPIMock(t, "testdata/import_flow")

	proc := startService(t, mock.server.URL+"/api/v2")

	t.Cleanup(func() { truncateTables(t) })

	// when: POST /imports starts a pokeapi import
	resp := doPost(t, proc.URL()+"/imports", `{"source": "pokeapi"}`)
	testastic.Equal(t, http.StatusCreated, resp.StatusCode)
	body := readBody(t, resp)
	testastic.AssertJSON(t, "testdata/import_flow/create_import_response.json", body)

	var importResp createdImportResponse

	decodeJSON(t, body, &importResp)
	assertUUIDV7(t, importResp.ID)
	testastic.Equal(t, "/imports/"+importResp.ID, resp.Header.Get("Location"))

	// then: the import completes and imported pokemon are available through the API
	importID := importResp.ID

	testastic.EventuallyEqual(t, "completed", func() string {
		resp := doGet(t, proc.URL()+"/imports/"+importID)
		body := readBody(t, resp)

		var importResp importStatusResponse

		decodeJSON(t, body, &importResp)

		if importResp.Status == "failed" {
			t.Fatalf("import %s failed", importID)
		}

		return importResp.Status
	}, 30*time.Second)

	resp = doGet(t, proc.URL()+"/imports/"+importID)
	testastic.Equal(t, http.StatusOK, resp.StatusCode)
	testastic.AssertJSON(t, "testdata/import_flow/completed_import_response.json", readBody(t, resp))

	resp = doGet(t, proc.URL()+"/pokemon")
	testastic.Equal(t, http.StatusOK, resp.StatusCode)
	testastic.AssertJSON(t, "testdata/import_flow/list_pokemon_response.json", readBody(t, resp))

	resp = doGet(t, proc.URL()+"/pokemon/1")
	testastic.Equal(t, http.StatusOK, resp.StatusCode)
	testastic.AssertJSON(t, "testdata/import_flow/get_bulbasaur_response.json", readBody(t, resp))
}

func TestListPokemonEmpty(t *testing.T) {
	// given: a running service with an empty database
	mock := newPokeAPIMock(t)
	proc := startService(t, mock.server.URL+"/api/v2")

	t.Cleanup(func() { truncateTables(t) })

	// when: GET /pokemon is called
	resp := doGet(t, proc.URL()+"/pokemon")

	// then: the pokemon list is empty
	testastic.Equal(t, http.StatusOK, resp.StatusCode)
	testastic.AssertJSON(t, "testdata/list_pokemon_empty/response.json", readBody(t, resp))
}

func TestListPokemonPaginated(t *testing.T) {
	// given: a running service with imported pokemon
	mock := newScenarioPokeAPIMock(t, "testdata/list_pokemon_paginated")
	proc := startService(t, mock.server.URL+"/api/v2")

	t.Cleanup(func() { truncateTables(t) })
	importPokemonForSetup(t, proc.URL())

	// when: GET /pokemon is called with a limit
	resp := doGet(t, proc.URL()+"/pokemon?limit=1")

	// then: the API returns the first page with the correct total
	testastic.Equal(t, http.StatusOK, resp.StatusCode)
	testastic.AssertJSON(t, "testdata/list_pokemon_paginated/response.json", readBody(t, resp))
}

func TestListPokemonFilteredByRarity(t *testing.T) {
	// given: a running service with imported pokemon of different rarities
	mock := newScenarioPokeAPIMock(t, "testdata/list_pokemon_filtered_by_rarity")
	proc := startService(t, mock.server.URL+"/api/v2")

	t.Cleanup(func() { truncateTables(t) })
	importPokemonForSetup(t, proc.URL())

	// when: GET /pokemon is filtered by common rarity
	resp := doGet(t, proc.URL()+"/pokemon?rarity=common")

	// then: the API returns only matching pokemon
	testastic.Equal(t, http.StatusOK, resp.StatusCode)
	testastic.AssertJSON(t, "testdata/list_pokemon_filtered_by_rarity/response.json", readBody(t, resp))
}

func TestGetPokemonNotFound(t *testing.T) {
	// given: a running service with no matching pokemon
	mock := newPokeAPIMock(t)
	proc := startService(t, mock.server.URL+"/api/v2")

	// when: GET /pokemon/9999 is called
	resp := doGet(t, proc.URL()+"/pokemon/9999")

	// then: the API returns a not found problem response
	testastic.Equal(t, http.StatusNotFound, resp.StatusCode)
	testastic.AssertJSON(t, "testdata/get_pokemon_not_found/response.json", readBody(t, resp))
}

func TestCreateCatchNoPokemon(t *testing.T) {
	// given: a running service with no imported pokemon
	mock := newPokeAPIMock(t)
	proc := startService(t, mock.server.URL+"/api/v2")

	t.Cleanup(func() { truncateTables(t) })

	// when: POST /catches is called
	resp := doPost(t, proc.URL()+"/catches", `{"pokeball_type": "pokeball"}`)

	// then: the API reports that no pokemon are available to catch
	testastic.Equal(t, http.StatusConflict, resp.StatusCode)
	testastic.AssertJSON(t, "testdata/create_catch_no_pokemon/response.json", readBody(t, resp))
}

func TestCreateCatchInvalidBody(t *testing.T) {
	// given: a running service
	mock := newPokeAPIMock(t)
	proc := startService(t, mock.server.URL+"/api/v2")

	// when: POST /catches is called with an empty body
	resp := doPost(t, proc.URL()+"/catches", "")

	// then: the API returns a bad request problem response
	testastic.Equal(t, http.StatusBadRequest, resp.StatusCode)
	testastic.AssertJSON(t, "testdata/create_catch_invalid_body/response.json", readBody(t, resp))
}

func TestCreateCatchAfterImport(t *testing.T) {
	// given: a running service and a PokeAPI fake serving imported pokemon
	mock := newCatchAfterImportMock(t)

	proc := startService(t, mock.server.URL+"/api/v2")

	t.Cleanup(func() { truncateTables(t) })
	importPokemonForSetup(t, proc.URL())

	// when: POST /catches is called after the import
	resp := doPost(t, proc.URL()+"/catches", `{"pokeball_type": "pokeball"}`)

	// then: the API creates a persisted catch and returns it
	testastic.Equal(t, http.StatusCreated, resp.StatusCode)
	body := readBody(t, resp)
	testastic.AssertJSON(t, "testdata/create_catch_after_import/create_response.json", body)

	var catchResp createdCatchResponse

	decodeJSON(t, body, &catchResp)
	assertUUIDV7(t, catchResp.ID)
	testastic.Equal(t, "/catches/"+catchResp.ID, resp.Header.Get("Location"))

	getResp := doGet(t, proc.URL()+"/catches/"+catchResp.ID)
	testastic.Equal(t, http.StatusOK, getResp.StatusCode)
	testastic.AssertJSON(t, "testdata/get_catch/existing/response.json", readBody(t, getResp))
}

func TestCreateImportInvalidBody(t *testing.T) {
	// given: a running service
	mock := newPokeAPIMock(t)
	proc := startService(t, mock.server.URL+"/api/v2")

	// when: POST /imports is called with a missing source
	resp := doPost(t, proc.URL()+"/imports", `{}`)

	// then: the API returns a bad request problem response
	testastic.Equal(t, http.StatusBadRequest, resp.StatusCode)
	testastic.AssertJSON(t, "testdata/import_invalid_body/response.json", readBody(t, resp))
}

func TestCreateImportUnsupportedSource(t *testing.T) {
	// given: a running service
	mock := newPokeAPIMock(t)
	proc := startService(t, mock.server.URL+"/api/v2")

	// when: POST /imports is called with an unsupported source
	resp := doPost(t, proc.URL()+"/imports", `{"source": "manual"}`)

	// then: the API rejects the unsupported source value
	testastic.Equal(t, http.StatusBadRequest, resp.StatusCode)
	testastic.AssertJSON(t, "testdata/import_unsupported_source/response.json", readBody(t, resp))
}

func TestGetImportNotFound(t *testing.T) {
	// given: a running service with no matching import
	mock := newPokeAPIMock(t)
	proc := startService(t, mock.server.URL+"/api/v2")

	// when: GET /imports/{id} is called for a missing import
	resp := doGet(t, proc.URL()+"/imports/550e8400-e29b-41d4-a716-446655440000")

	// then: the API returns a not found problem response
	testastic.Equal(t, http.StatusNotFound, resp.StatusCode)
	testastic.AssertJSON(t, "testdata/get_import_not_found/response.json", readBody(t, resp))
}

func TestGetCatchNotFound(t *testing.T) {
	// given: a running service with no matching catch
	mock := newPokeAPIMock(t)
	proc := startService(t, mock.server.URL+"/api/v2")

	// when: GET /catches/{id} is called for a missing catch
	resp := doGet(t, proc.URL()+"/catches/550e8400-e29b-41d4-a716-446655440000")

	// then: the API returns a not found problem response
	testastic.Equal(t, http.StatusNotFound, resp.StatusCode)
	testastic.AssertJSON(t, "testdata/get_catch/not_found/response.json", readBody(t, resp))
}

// Minimal response types for workflow values needed across requests.

type createdImportResponse struct {
	ID string `json:"id"`
}

type createdCatchResponse struct {
	ID string `json:"id"`
}

type importStatusResponse struct {
	Status string `json:"status"`
}

func assertUUIDV7(t *testing.T, raw string) {
	t.Helper()

	id, err := uuid.Parse(raw)
	testastic.NoError(t, err)
	testastic.Equal(t, uuid.Version(7), id.Version())
}

// HTTP helpers.

func doGet(t *testing.T, url string) *http.Response {
	t.Helper()

	resp, err := http.Get(url) //nolint:noctx // Test code.
	testastic.NoError(t, err)

	t.Cleanup(func() { resp.Body.Close() })

	return resp
}

func doPost(t *testing.T, url string, body string) *http.Response {
	t.Helper()

	resp, err := http.Post(url, "application/json", strings.NewReader(body)) //nolint:noctx // Test code.
	testastic.NoError(t, err)

	t.Cleanup(func() { resp.Body.Close() })

	return resp
}

func readBody(t *testing.T, resp *http.Response) []byte {
	t.Helper()

	body, err := io.ReadAll(resp.Body)
	testastic.NoError(t, err)

	return body
}

func decodeJSON(t *testing.T, body []byte, v any) {
	t.Helper()

	err := json.Unmarshal(body, v)
	testastic.NoError(t, err)
}
