package repository

import (
	"context"
	"database/sql"
	"errors"
	"goauth/model"
	"time"

	"goauth/db"
)

type UserRepository struct {
	q db.Querier
}

type TokenRepository struct {
	q db.Querier
}

func NewUserRepository(dbConn *sql.DB) *UserRepository {
	return &UserRepository{q: db.New(dbConn)}
}

func NewTokenRepository(dbConn *sql.DB) *TokenRepository {
	return &TokenRepository{q: db.New(dbConn)}
}

func (r *UserRepository) GetOrCreateUser(ctx context.Context, email, firstName, lastName string) (model.User, error) {
	user, err := r.q.GetUserByEmail(ctx, email)
	if errors.Is(err, sql.ErrNoRows) {
		user, err = r.q.CreateUser(ctx, db.CreateUserParams{
			Email:     email,
			FirstName: firstName,
			LastName:  lastName,
		})
	}

	return model.User{
		ID:        user.ID,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
	}, err
}

func (r *UserRepository) GetUserByID(ctx context.Context, id string) (model.User, error) {
	user, err := r.q.GetUser(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return model.User{}, err
	}

	return model.User{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}, err
}

func (r *TokenRepository) StoreRefreshToken(ctx context.Context, userID string, refreshToken string) error {
	_, err := r.q.CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
		UserID:    userID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	})
	return err
}

func (r *TokenRepository) VerifyRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	token, err := r.q.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", err
	}
	return token.UserID, nil
}

func (r *TokenRepository) UpdateRefreshToken(ctx context.Context, oldRefreshToken, newRefreshToken string) error {
	_, err := r.q.UpdateRefreshToken(ctx, db.UpdateRefreshTokenParams{
		Token:     oldRefreshToken,
		Token_2:   newRefreshToken,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	})
	return err
}
