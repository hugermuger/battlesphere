-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password, user_name)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetUserByUserName :one
SELECT * FROM users
WHERE user_name = $1;
