package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                         uuid.UUID  `json:"id"`
	Email                      string     `json:"email"`
	PasswordHash               string     `json:"-"`
	IsVerified                 bool       `json:"is_verified"`
	EmailVerificationCode      *string    `json:"-"`
	EmailVerificationExpiresAt *time.Time `json:"-"`
	CreatedAt                  time.Time  `json:"created_at"`
	UpdatedAt                  time.Time  `json:"updated_at"`
}
