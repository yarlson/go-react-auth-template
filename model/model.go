package model

import (
	"time"
)

type User struct {
	ID string

	Email     string
	FirstName string
	LastName  string
}

type RefreshToken struct {
	UserID    uint
	Token     string
	ExpiresAt time.Time
}
