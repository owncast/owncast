package utils

import (
	"crypto/sha512"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	password_bytes := compressPassword([]byte(password))
	// 0 will use the default cost of 10 instead
	hash, err := bcrypt.GenerateFromPassword(password_bytes, 0)
	return string(hash), err
}

func CompareHash(hash string, password string) error {
	password_bytes := compressPassword([]byte(password))
	return bcrypt.CompareHashAndPassword([]byte(hash), password_bytes)
}

// Takes a password and computes a sha-512 hash of it if it is longer than 72 bytes, guaranteeing it is less than 72 bytes long.
func compressPassword(password []byte) []byte {
	if len(password) > 72 {
		sha512_hashed := sha512.Sum512(password)
		password = sha512_hashed[:]
	}
	return password
}
