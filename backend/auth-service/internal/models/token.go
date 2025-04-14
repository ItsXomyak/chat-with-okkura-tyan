package models

import (
	"time"

	"github.com/google/uuid"
)

type ConfirmationToken struct {
	ID        int64     `json:"id" db:"id"`
	UserID    uuid.UUID `json:"userId" db:"user_id"`
	Token     string    `json:"token" db:"token"`
	ExpiresAt time.Time `json:"expiresAt" db:"expires_at"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

type PasswordResetToken struct {
	ID        int64     `json:"id" db:"id"`
	UserID    uuid.UUID `json:"userId" db:"user_id"`
	Token     string    `json:"token" db:"token"`
	ExpiresAt time.Time `json:"expiresAt" db:"expires_at"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	Used      bool      `json:"used" db:"used"`
}
