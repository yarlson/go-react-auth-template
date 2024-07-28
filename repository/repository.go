package repository

import (
	"fmt"
	"goauth/model"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Email     string `gorm:"uniqueIndex"`
	FirstName string
	LastName  string
}

type RefreshToken struct {
	gorm.Model

	UserID    string
	Token     string `gorm:"uniqueIndex"`
	ExpiresAt time.Time
}

type UserRepository struct {
	db *gorm.DB
}

type TokenRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func NewTokenRepository(db *gorm.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *UserRepository) GetOrCreateUser(email, firstName, lastName string) (model.User, error) {
	var user User
	err := r.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Where(User{Email: email}).FirstOrCreate(&user, User{Email: email, FirstName: firstName, LastName: lastName})
		return result.Error
	})

	return model.User{
		ID:        fmt.Sprintf("%d", user.ID),
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
	}, err
}

func (r *UserRepository) GetUserByID(id string) (model.User, error) {
	var user User
	result := r.db.First(&user, id)

	return model.User{
		ID:        fmt.Sprintf("%d", user.ID),
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}, result.Error
}

func (r *TokenRepository) StoreRefreshToken(userID string, refreshToken string) error {
	token := RefreshToken{
		UserID:    userID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour), // 30 days expiration
	}

	return r.db.Create(&token).Error
}

func (r *TokenRepository) VerifyRefreshToken(refreshToken string) (string, error) {
	var token RefreshToken
	result := r.db.Where("token = ? AND expires_at > ?", refreshToken, time.Now()).First(&token)
	if result.Error != nil {
		return "", result.Error
	}

	return token.UserID, nil
}

func (r *TokenRepository) UpdateRefreshToken(userID string, oldRefreshToken, newRefreshToken string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("token = ?", oldRefreshToken).Delete(&RefreshToken{}).Error; err != nil {
			return err
		}

		token := RefreshToken{
			UserID:    userID,
			Token:     newRefreshToken,
			ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
		}

		return tx.Create(&token).Error
	})
}
