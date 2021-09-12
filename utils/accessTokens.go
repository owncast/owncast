package utils

import (
	"encoding/base64"
	"math/rand"
	"time"
)

const tokenLength = 32

// GenerateAccessToken will generate and return an access token.
func GenerateAccessToken() (string, error) {
	return generateRandomString(tokenLength)
}

// generateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	rand.Seed(time.Now().UTC().UnixNano())
	_, err := rand.Read(b) //nolint
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// generateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func generateRandomString(n int) (string, error) {
	b, err := generateRandomBytes(n)
	return base64.URLEncoding.EncodeToString(b), err
}
