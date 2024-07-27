package repository

import (
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

	UserID    uint
	Token     string `gorm:"uniqueIndex"`
	ExpiresAt time.Time
}

type UserRepository interface {
	GetOrCreateUser(email, firstName, lastName string) (User, error)
	GetUserByID(id uint) (User, error)
}

type TokenRepository interface {
	StoreRefreshToken(userID uint, refreshToken string) error
	VerifyRefreshToken(refreshToken string) (uint, error)
	UpdateRefreshToken(userID uint, newRefreshToken string) error
}

type GormUserRepository struct {
	db *gorm.DB
}

type GormTokenRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

func NewGormTokenRepository(db *gorm.DB) *GormTokenRepository {
	return &GormTokenRepository{db: db}
}

func (r *GormUserRepository) GetOrCreateUser(email, firstName, lastName string) (User, error) {
	var user User
	err := r.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Where(User{Email: email}).FirstOrCreate(&user, User{Email: email, FirstName: firstName, LastName: lastName})
		return result.Error
	})

	return user, err
}

func (r *GormUserRepository) GetUserByID(id uint) (User, error) {
	var user User
	result := r.db.First(&user, id)

	return user, result.Error
}

func (r *GormTokenRepository) StoreRefreshToken(userID uint, refreshToken string) error {
	token := RefreshToken{
		UserID:    userID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour), // 30 days expiration
	}

	return r.db.Create(&token).Error
}

func (r *GormTokenRepository) VerifyRefreshToken(refreshToken string) (uint, error) {
	var token RefreshToken
	result := r.db.Where("token = ? AND expires_at > ?", refreshToken, time.Now()).First(&token)
	if result.Error != nil {
		return 0, result.Error
	}

	return token.UserID, nil
}

func (r *GormTokenRepository) UpdateRefreshToken(userID uint, oldRefreshToken, newRefreshToken string) error {
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
