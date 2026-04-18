-- +goose Up
CREATE TABLE pokemon (
    pokedex_id      INTEGER PRIMARY KEY,
    name            TEXT NOT NULL UNIQUE,
    rarity          TEXT NOT NULL CHECK (rarity IN ('common', 'uncommon', 'rare', 'legendary', 'mythical')),
    types           TEXT[] NOT NULL,
    sprite_url      TEXT NOT NULL DEFAULT '',
    hp              INTEGER NOT NULL,
    attack          INTEGER NOT NULL,
    defense         INTEGER NOT NULL,
    special_attack  INTEGER NOT NULL,
    special_defense INTEGER NOT NULL,
    speed           INTEGER NOT NULL,
    base_experience INTEGER NOT NULL DEFAULT 0,
    capture_rate    INTEGER NOT NULL DEFAULT 0,
    is_legendary    BOOLEAN NOT NULL DEFAULT FALSE,
    is_mythical     BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_pokemon_rarity ON pokemon (rarity);

-- +goose Down
DROP TABLE IF EXISTS pokemon;
