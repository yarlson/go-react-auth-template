-- name: GetUser :one
SELECT *
FROM users
WHERE id = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (id, email, first_name, last_name)
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT *
FROM users
ORDER BY id;

-- name: UpdateUser :exec
UPDATE users
SET first_name = $2,
    last_name  = $3
WHERE id = $1;

-- name: DeleteUser :exec
DELETE
FROM users
WHERE id = $1;