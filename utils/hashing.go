package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	// 0 will use the default cost of 10 instead
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	return string(hash), err
}

func ComparseHash(hash string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
