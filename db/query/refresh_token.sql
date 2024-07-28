-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (user_id, token, expires_at)
VALUES ($1, $2, $3) RETURNING *;

-- name: GetRefreshToken :one
SELECT *
FROM refresh_tokens
WHERE token = $1
  AND expires_at > NOW() LIMIT 1;

-- name: DeleteRefreshToken :exec
DELETE
FROM refresh_tokens
WHERE token = $1;

-- name: UpdateRefreshToken :one
UPDATE refresh_tokens
SET token      = $2,
    expires_at = $3
WHERE token = $1 RETURNING *;

-- name: DeleteExpiredTokens :exec
DELETE
FROM refresh_tokens
WHERE expires_at <= NOW();