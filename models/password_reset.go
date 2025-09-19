package models

import (
	"database/sql"
	"fmt"
	"time"
)

const (
	DefaultResetDuration = time.Hour
)

type PasswordReset struct {
	ID     int
	UserID int
	//Token is only set when a password reset happens
	Token     string
	TokenHash string
	ExpiresAt time.Time
}

type PasswordResetService struct {
	DB            *sql.DB
	BytesPerToken int
	// Amount of time that a password reset is valid for
	// defaults to DefaultResetDuration
	Duration time.Duration
}

func (pr *PasswordResetService) Create(email string) (*PasswordReset, error) {
	return nil, fmt.Errorf("TODO: Implement password reset service.Create")
}

func (pr *PasswordResetService) Consume(token string) (*User, error) {
	return nil, fmt.Errorf("TODO: Implement password reset service.Consume")
}
