-- name: UpsertPokemon :exec
INSERT INTO pokemon (
    pokedex_id, name, rarity, types, sprite_url,
    hp, attack, defense, special_attack, special_defense, speed,
    base_experience, capture_rate, is_legendary, is_mythical
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
ON CONFLICT (pokedex_id) DO UPDATE SET
    name = EXCLUDED.name,
    rarity = EXCLUDED.rarity,
    types = EXCLUDED.types,
    sprite_url = EXCLUDED.sprite_url,
    hp = EXCLUDED.hp,
    attack = EXCLUDED.attack,
    defense = EXCLUDED.defense,
    special_attack = EXCLUDED.special_attack,
    special_defense = EXCLUDED.special_defense,
    speed = EXCLUDED.speed,
    base_experience = EXCLUDED.base_experience,
    capture_rate = EXCLUDED.capture_rate,
    is_legendary = EXCLUDED.is_legendary,
    is_mythical = EXCLUDED.is_mythical,
    updated_at = NOW();

-- name: GetPokemonByID :one
SELECT pokedex_id, name, rarity, types, sprite_url,
    hp, attack, defense, special_attack, special_defense, speed,
    base_experience, capture_rate, is_legendary, is_mythical,
    created_at, updated_at
FROM pokemon
WHERE pokedex_id = $1;

-- name: ListPokemon :many
SELECT pokedex_id, name, rarity, types, sprite_url,
    hp, attack, defense, special_attack, special_defense, speed,
    base_experience, capture_rate, is_legendary, is_mythical,
    created_at, updated_at
FROM pokemon
ORDER BY pokedex_id
LIMIT $1 OFFSET $2;

-- name: ListPokemonByRarity :many
SELECT pokedex_id, name, rarity, types, sprite_url,
    hp, attack, defense, special_attack, special_defense, speed,
    base_experience, capture_rate, is_legendary, is_mythical,
    created_at, updated_at
FROM pokemon
WHERE rarity = $1
ORDER BY pokedex_id
LIMIT $2 OFFSET $3;

-- name: CountPokemon :one
SELECT COUNT(*) FROM pokemon;

-- name: CountPokemonByRarity :one
SELECT COUNT(*) FROM pokemon WHERE rarity = $1;

-- name: GetRandomPokemonByRarity :one
SELECT pokedex_id, name, rarity, types, sprite_url,
    hp, attack, defense, special_attack, special_defense, speed,
    base_experience, capture_rate, is_legendary, is_mythical,
    created_at, updated_at
FROM pokemon
WHERE rarity = $1
ORDER BY RANDOM()
LIMIT 1;

-- name: CreateImport :exec
INSERT INTO imports (id, source, status, item_count, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetImport :one
SELECT id, source, status, item_count, created_at, updated_at
FROM imports
WHERE id = $1;

-- name: UpdateImportStatus :exec
UPDATE imports
SET status = $2, item_count = $3, updated_at = NOW()
WHERE id = $1;
