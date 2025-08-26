package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"

	"github.com/rahulbalajee/lenslocked/rand"
)

const (
	// Min number of bytes per each session token
	MinBytesPerToken = 32
)

type Session struct {
	ID     int
	UserID int
	// Token is only set when creating a new session, otherwise it will be empty
	Token     string
	TokenHash string
}

type SessionService struct {
	DB            *sql.DB
	BytesPerToken int
}

func (ss *SessionService) Create(userID int) (*Session, error) {
	bytesPerToken := ss.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}

	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	session := Session{
		UserID:    userID,
		Token:     token,
		TokenHash: ss.hash(token),
	}

	// 1) Query for a user's session
	// If found update it
	// If not found create a new session

	// Or skip (1) and try update if error create session
	row := ss.DB.QueryRow(`
		UPDATE sessions
		SET token_hash = $2
		WHERE user_id = $1
		RETURNING id;`, session.UserID, session.TokenHash)
	err = row.Scan(&session.ID)
	if err == sql.ErrNoRows {
		row := ss.DB.QueryRow(`
			INSERT INTO sessions (user_id, token_hash)
			VALUES ($1, $2)
			RETURNING id;`, session.UserID, session.TokenHash)
		err = row.Scan(&session.ID)
	}
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	return &session, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	// Hash the session token
	tokenHash := ss.hash(token)

	// Query for that session wth the hash
	var user User

	row := ss.DB.QueryRow(`SELECT user_id FROM sessions WHERE token_hash = $1;`, tokenHash)

	err := row.Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}

	// Using the UserID from the session query for that user
	row = ss.DB.QueryRow(`
		SELECT email, password_hash
		FROM users where id = $1`, user.ID)

	err = row.Scan(&user.Email, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}

	// Return the user
	return &user, nil
}

func (ss *SessionService) Delete(token string) error {
	tokenHash := ss.hash(token)

	_, err := ss.DB.Exec(`
		DELETE FROM sessions
		WHERE token_hash = $1;`, tokenHash)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

func (ss *SessionService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
