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

func (tm *TokenManager) New() (token, tokenHash string, err error) {
	bytesPerToken := tm.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}

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
