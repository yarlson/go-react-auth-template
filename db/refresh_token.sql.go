// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: refresh_token.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createRefreshToken = `-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (id, user_id, token, expires_at)
VALUES ($1, $2, $3, $4) RETURNING id, user_id, token, expires_at
`

type CreateRefreshTokenParams struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (q *Queries) CreateRefreshToken(ctx context.Context, arg CreateRefreshTokenParams) (RefreshToken, error) {
	row := q.db.QueryRowContext(ctx, createRefreshToken,
		arg.ID,
		arg.UserID,
		arg.Token,
		arg.ExpiresAt,
	)
	var i RefreshToken
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Token,
		&i.ExpiresAt,
	)
	return i, err
}

const deleteExpiredTokens = `-- name: DeleteExpiredTokens :exec
DELETE
FROM refresh_tokens
WHERE expires_at <= NOW()
`

func (q *Queries) DeleteExpiredTokens(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteExpiredTokens)
	return err
}

const deleteRefreshToken = `-- name: DeleteRefreshToken :exec
DELETE
FROM refresh_tokens
WHERE token = $1
`

func (q *Queries) DeleteRefreshToken(ctx context.Context, token string) error {
	_, err := q.db.ExecContext(ctx, deleteRefreshToken, token)
	return err
}

const getRefreshToken = `-- name: GetRefreshToken :one
SELECT id, user_id, token, expires_at
FROM refresh_tokens
WHERE token = $1
  AND expires_at > NOW() LIMIT 1
`

func (q *Queries) GetRefreshToken(ctx context.Context, token string) (RefreshToken, error) {
	row := q.db.QueryRowContext(ctx, getRefreshToken, token)
	var i RefreshToken
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Token,
		&i.ExpiresAt,
	)
	return i, err
}

const updateRefreshToken = `-- name: UpdateRefreshToken :one
UPDATE refresh_tokens
SET token      = $2,
    expires_at = $3
WHERE token = $1 RETURNING id, user_id, token, expires_at
`

type UpdateRefreshTokenParams struct {
	Token     string    `json:"token"`
	Token_2   string    `json:"token_2"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (q *Queries) UpdateRefreshToken(ctx context.Context, arg UpdateRefreshTokenParams) (RefreshToken, error) {
	row := q.db.QueryRowContext(ctx, updateRefreshToken, arg.Token, arg.Token_2, arg.ExpiresAt)
	var i RefreshToken
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Token,
		&i.ExpiresAt,
	)
	return i, err
}
