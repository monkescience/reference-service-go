-- +goose Up
CREATE TABLE catches (
    id                UUID PRIMARY KEY,
    pokemon_pokedex_id INTEGER NOT NULL REFERENCES pokemon (pokedex_id),
    pokeball_type     TEXT NOT NULL CHECK (pokeball_type IN ('pokeball', 'great_ball', 'ultra_ball', 'master_ball')),
    is_shiny          BOOLEAN NOT NULL,
    caught_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS catches;
