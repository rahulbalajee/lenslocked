package rand

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)

	nRead, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("rand bytes: %w", err)
	}
	if nRead < n {
		return nil, fmt.Errorf("rand bytes: didn't read enough random bytes")
	}

	return b, nil
}

// n is the number of bytes used to generate rand string/bytes
// String() converts the random bytes generated into a string and returns to TokenManager
func String(n int) (string, error) {
	b, err := Bytes(n)
	if err != nil {
		return "", fmt.Errorf("rand string: %w", err)
	}

	return base64.URLEncoding.EncodeToString(b), nil
}
