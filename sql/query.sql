-- name: ListUsers :many
SELECT * FROM users;

-- name: CreateUser :one
INSERT INTO users (id, email, username, profileIndex)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ?;

-- name: UpdateUser :exec
UPDATE users
SET email = ?, username = ?, profileIndex = ?
WHERE email = ?;
