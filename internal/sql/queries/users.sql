-- name: CreateUser :one
INSERT INTO users(id, login, hashed_password, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetUserByLogin :one
SELECT * FROM users
WHERE login = $1;