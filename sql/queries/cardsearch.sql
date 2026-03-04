-- name: SearchCardsByNameListEng :many
SELECT name, set_name, set_code, rarity FROM cards
WHERE name ILIKE '%' || $1 || '%' AND lang = 'en' AND 'paper' = ANY(games)
ORDER BY release_date DESC
LIMIT $2 OFFSET $3;

-- name: CountCardsByNameListEng :one
SELECT COUNT(*) FROM cards
WHERE name ILIKE '%' || $1 || '%' AND lang = 'en' AND 'paper' = ANY(games);

-- name: SearchCardsByNameList :many
SELECT printed_name, set_name, set_code, rarity FROM cards
WHERE printed_name ILIKE '%' || $1 || '%' AND lang = $2 AND 'paper' = ANY(games)
ORDER BY release_date DESC
LIMIT $3 OFFSET $4;

-- name: CountCardsByNameList :one
SELECT COUNT(*)  FROM cards
WHERE printed_name ILIKE '%' || $1 || '%' AND lang = $2 AND 'paper' = ANY(games);

-- name: DoesLangExist :one
SELECT name FROM cards WHERE lang = $1 LIMIT 1;
