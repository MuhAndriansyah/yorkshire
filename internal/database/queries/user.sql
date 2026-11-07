-- name: CreateUser :one
INSERT INTO users (
  id, username, email, password_hash
) VALUES (
  $1, $2, $3, $4
)
RETURNING id, username, email, password_hash, avatar_url, created_at, updated_at;

-- name: GetUserByEmail :one
SELECT id, username, email, password_hash, avatar_url, created_at, updated_at
FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT id, username, email, avatar_url, created_at, updated_at
FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserEmailByID :one
SELECT email
FROM users
WHERE id = $1 LIMIT 1;
