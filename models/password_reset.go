package models

import (
	"database/sql"
	"fmt"
	"strings"
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
	DB *sql.DB
	// Amount of time that a password reset is valid for
	// defaults to DefaultResetDuration
	Duration     time.Duration
	TokenManager TokenManager
}

func (pr *PasswordResetService) Create(email string) (*PasswordReset, error) {
	email = strings.ToLower(email)

	var userID int

	row := pr.DB.QueryRow(`
		SELECT id FROM users WHERE email = $1`, email)

	err := row.Scan(&userID)
	if err != nil {
		// TODO: consider returning a different error when user doesn't exists
		return nil, fmt.Errorf("create password reset: %w", err)
	}

	token, tokenHash, err := pr.TokenManager.New()
	if err != nil {
		return nil, fmt.Errorf("create password reset: %w", err)
	}

	duration := pr.Duration
	if duration == 0 {
		duration = DefaultResetDuration
	}

	pwReset := PasswordReset{
		UserID:    userID,
		Token:     token,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(duration),
	}

	// Insert pwReset into DB
	row = pr.DB.QueryRow(`
		INSERT INTO password_resets (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3) ON CONFLICT (user_id) DO
		UPDATE
		SET token_hash = $2, expires_at = $3
		RETURNING id`, pwReset.UserID, pwReset.TokenHash, pwReset.ExpiresAt)

	err = row.Scan(&pwReset.ID)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	return &pwReset, nil
}

func (pr *PasswordResetService) Consume(token string) (*User, error) {
	tokenHash := pr.TokenManager.Hash(token)

	var user User
	var pwReset PasswordReset

	row := pr.DB.QueryRow(`
		SELECT password_resets.id,
			password_resets.expires_at,
			users.id,
			users.email,
			users.password_hash
		FROM password_resets 
			JOIN users on users.id = password_resets.user_id
		WHERE password_resets.token_hash = $1;`, tokenHash)

	err := row.Scan(
		&pwReset.ID,
		&pwReset.ExpiresAt,
		&user.ID,
		&user.Email,
		&user.PasswordHash,
	)
	if err != nil {
		return nil, fmt.Errorf("consume password reset: %w", err)
	}

	if time.Now().After(pwReset.ExpiresAt) {
		return nil, fmt.Errorf("token expired: %v", token)
	}

	err = pr.delete(pwReset.ID)
	if err != nil {
		return nil, fmt.Errorf("consume password reset: %w", err)
	}

	return &user, nil
}

func (pr *PasswordResetService) delete(id int) error {
	_, err := pr.DB.Exec(`
		DELETE FROM password_resets
		WHERE id = $1`, id)

	if err != nil {
		return fmt.Errorf("delete password reset: %w", err)
	}

	return nil
}
