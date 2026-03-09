-- name: SearchCardsByNameListEng :many
SELECT * FROM (
    SELECT DISTINCT ON (oracle_id)
        id,
        oracle_id,
        name,
        mana_cost,
        type_line,
        release_date,
        layout
    FROM cards
    WHERE name ILIKE '%' || $1 || '%'
      AND lang = 'en'
      AND 'paper' = ANY(games)
      AND layout != 'art_series'
    ORDER BY oracle_id, release_date DESC
) AS unique_cards
ORDER BY name ASC
LIMIT $2 OFFSET $3;

-- name: CountCardsByNameListEng :one
SELECT COUNT(DISTINCT oracle_id)
FROM cards
WHERE name ILIKE '%' || $1 || '%'
    AND lang = 'en'
    AND 'paper' = ANY(games)
    AND layout != 'art_series';

-- name: SearchCardsByNameList :many
SELECT * FROM (
    SELECT DISTINCT ON (oracle_id)
        oracle_id,
        printed_name,
        mana_cost,
        printed_type_line,
        layout
    FROM cards
    WHERE printed_name ILIKE '%' || $1 || '%'
        AND lang = $2
        AND 'paper' = ANY(games)
        AND layout != 'art_series'
    ORDER BY oracle_id, release_date DESC
) AS unique_cards
ORDER BY printed_name ASC
LIMIT $3 OFFSET $4;

-- name: CountCardsByNameList :one
SELECT COUNT(DISTINCT oracle_id)
FROM cards
WHERE printed_name ILIKE '%' || $1 || '%'
    AND lang = $2
    AND 'paper' = ANY(games)
    AND layout != 'art_series';

-- name: DoesLangExist :one
SELECT name
FROM cards
WHERE lang = $1 LIMIT 1;

-- name: SearchCardsByOracleIDList :many
SELECT id,
    name,
    flavor_name,
    printed_name,
    release_date,
    set_name,
    set_code,
    collector_number
FROM cards
WHERE oracle_id = $1
    AND lang = $2
    AND 'paper' = ANY(games)
ORDER BY release_date DESC
LIMIT $3 OFFSET $4;

-- name: SearchCardByOracleID :one
SELECT *
FROM cards
WHERE oracle_id = $1
    AND lang = $2
    AND 'paper' = ANY(games)
ORDER BY release_date DESC
LIMIT 1;

-- name: GetCardFaces :many
SELECT *
FROM card_faces
WHERE card_id = $1;

-- name: CountCardsByOracleIDList :one
SELECT COUNT(*)
FROM cards
WHERE oracle_id = $1
    AND lang = $2
    AND 'paper' = ANY(games);

-- name: GetCardByID :one
SELECT *
FROM cards
WHERE id = $1;

-- name: GetCardLegalties :one
SELECT * FROM legalities WHERE card_id = $1;

-- name: GetOracleRulings :many
SELECT * FROM rulings WHERE oracle_id = $1;
