package models

import (
	"crypto/sha256"
	"encoding/base64"

	"github.com/rahulbalajee/lenslocked/rand"
)

const (
	// Min number of bytes per each session token
	MinBytesPerToken = 32
)

type TokenManager struct {
	BytesPerToken int
}

// Pointer receiver to make sure BytesPerToken is passed to the method when set in a different function and called
func (tm *TokenManager) New() (token, tokenHash string, err error) {
	bytesPerToken := tm.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}

	// rand.String() uses the Crypto rand package from stdlib to generate a unique random string
	token, err = rand.String(bytesPerToken)
	if err != nil {
		return "", "", err
	}

	tokenHash = tm.Hash(token)

	return token, tokenHash, nil
}

func (tm *TokenManager) Hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
