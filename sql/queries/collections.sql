-- name: AddCardToCollections :exec
INSERT INTO collections (id, user_id, folder_id, scryfall_id, purchase_date, purchase_price, finish, condition, created_at)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    NOW()
);
