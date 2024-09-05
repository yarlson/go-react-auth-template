package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	BaseModel

	Email     string `gorm:"uniqueIndex" json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type RefreshToken struct {
	BaseModel

	UserID    uuid.UUID `gorm:"type:uuid" json:"userId"`
	Token     string    `gorm:"uniqueIndex" json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type SessionToken struct {
	BaseModel

	UserID    uuid.UUID `gorm:"type:uuid" json:"userId"`
	Token     string    `gorm:"uniqueIndex" json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}
