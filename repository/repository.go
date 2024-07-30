package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"goauth/db"
	"goauth/model"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrTokenInvalid = errors.New("invalid refresh token")
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
	if err == nil {
		return modelFromDBUser(user), nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return model.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	// User doesn't exist, create a new one
	id, err := uuid.NewV7()
	if err != nil {
		return model.User{}, fmt.Errorf("failed to generate UUID: %w", err)
	}

	newUser, err := r.q.CreateUser(ctx, db.CreateUserParams{
		ID:        id,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
	})
	if err != nil {
		return model.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return modelFromDBUser(newUser), nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id string) (model.User, error) {
	uUserId, err := uuid.Parse(id)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to parse UUID: %w", err)
	}

	user, err := r.q.GetUser(ctx, uUserId)
	if errors.Is(err, sql.ErrNoRows) {
		return model.User{}, ErrUserNotFound
	}
	if err != nil {
		return model.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return modelFromDBUser(user), nil
}

func modelFromDBUser(dbUser db.User) model.User {
	return model.User{
		ID:        dbUser.ID.String(),
		Email:     dbUser.Email,
		FirstName: dbUser.FirstName,
		LastName:  dbUser.LastName,
	}
}

func (r *TokenRepository) StoreRefreshToken(ctx context.Context, userID string, refreshToken string) error {
	id, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("failed to generate UUID: %w", err)
	}

	uUserId, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("failed to parse user ID: %w", err)
	}

	_, err = r.q.CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
		ID:        id,
		UserID:    uUserId,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	})
	if err != nil {
		return fmt.Errorf("failed to store refresh token: %w", err)
	}

	return nil
}

func (r *TokenRepository) VerifyRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	token, err := r.q.GetRefreshToken(ctx, refreshToken)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrTokenInvalid
	}

	if err != nil {
		return "", fmt.Errorf("failed to verify refresh token: %w", err)
	}

	if token.ExpiresAt.Before(time.Now()) {
		return "", ErrTokenInvalid
	}

	return token.UserID.String(), nil
}

func (r *TokenRepository) UpdateRefreshToken(ctx context.Context, oldRefreshToken, newRefreshToken string) error {
	_, err := r.q.UpdateRefreshToken(ctx, db.UpdateRefreshTokenParams{
		Token:     oldRefreshToken,
		Token_2:   newRefreshToken,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrTokenInvalid
		}
		return fmt.Errorf("failed to update refresh token: %w", err)
	}

	return nil
}
