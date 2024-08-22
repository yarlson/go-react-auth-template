package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	BaseModel

	Email     string
	FirstName string
	LastName  string
}

type RefreshToken struct {
	BaseModel

	UserID    uuid.UUID
	Token     string
	ExpiresAt time.Time
}
