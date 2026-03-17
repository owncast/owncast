package utils

import (
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword_ShortPassword(t *testing.T) {
	password := "short"
	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("HashPassword failed for short password: %v", err)
	}

	if hash == "" {
		t.Error("HashPassword returned empty hash for short password")
	}

	if hash == password {
		t.Error("HashPassword returned unhashed password")
	}
}

func TestHashPassword_72BytePassword(t *testing.T) {
	// Test with exactly 72 bytes (bcrypt's limit)
	password := strings.Repeat("a", 72)
	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("HashPassword failed for 72-byte password: %v", err)
	}

	if hash == "" {
		t.Error("HashPassword returned empty hash for 72-byte password")
	}

	if hash == password {
		t.Error("HashPassword returned unhashed password")
	}
}

func TestHashPassword_LongPassword(t *testing.T) {
	// Test with password longer than 72 bytes
	password := strings.Repeat("a", 100)
	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("HashPassword failed for long password: %v", err)
	}

	if hash == "" {
		t.Error("HashPassword returned empty hash for long password")
	}

	if hash == password {
		t.Error("HashPassword returned unhashed password")
	}
}

func TestHashPassword_VeryLongPassword(t *testing.T) {
	// Test with very long password (200 bytes)
	password := strings.Repeat("x", 200)
	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("HashPassword failed for very long password: %v", err)
	}

	if hash == "" {
		t.Error("HashPassword returned empty hash for very long password")
	}
}

func TestCompareHash_ShortPassword_Success(t *testing.T) {
	password := "testpassword123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	err = CompareHash(hash, password)
	if err != nil {
		t.Errorf("CompareHash failed to match correct short password: %v", err)
	}
}

func TestCompareHash_ShortPassword_Failure(t *testing.T) {
	password := "testpassword123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	err = CompareHash(hash, "wrongpassword")
	if err == nil {
		t.Error("CompareHash incorrectly matched wrong password")
	}
}

func TestCompareHash_LongPassword_Success(t *testing.T) {
	// Test with password longer than 72 bytes
	password := strings.Repeat("a", 100)
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed for long password: %v", err)
	}

	err = CompareHash(hash, password)
	if err != nil {
		t.Errorf("CompareHash failed to match correct long password: %v", err)
	}
}

func TestCompareHash_LongPassword_Failure(t *testing.T) {
	// Test with password longer than 72 bytes
	password := strings.Repeat("a", 100)
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	// Wrong password also longer than 72 bytes
	wrongPassword := strings.Repeat("b", 100)
	err = CompareHash(hash, wrongPassword)
	if err == nil {
		t.Error("CompareHash incorrectly matched wrong long password")
	}
}

func TestCompareHash_LongPassword_SlightDifference(t *testing.T) {
	// Test that even a slight difference in long passwords is detected
	password := strings.Repeat("a", 100)
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	// Same length but one character different
	wrongPassword := strings.Repeat("a", 99) + "b"
	err = CompareHash(hash, wrongPassword)
	if err == nil {
		t.Error("CompareHash incorrectly matched password with slight difference")
	}
}

func TestCompareHash_73BytePassword(t *testing.T) {
	// Test edge case: 73 bytes (just over bcrypt's limit)
	password := strings.Repeat("a", 73)
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed for 73-byte password: %v", err)
	}

	err = CompareHash(hash, password)
	if err != nil {
		t.Errorf("CompareHash failed to match correct 73-byte password: %v", err)
	}

	// Ensure wrong password doesn't match
	wrongPassword := strings.Repeat("b", 73)
	err = CompareHash(hash, wrongPassword)
	if err == nil {
		t.Error("CompareHash incorrectly matched wrong 73-byte password")
	}
}

func TestCompressPassword_ShortPassword(t *testing.T) {
	password := []byte("short")
	compressed := compressPassword(password)

	if len(compressed) != len(password) {
		t.Error("compressPassword should not modify passwords shorter than 72 bytes")
	}

	if string(compressed) != string(password) {
		t.Error("compressPassword modified short password content")
	}
}

func TestCompressPassword_72BytePassword(t *testing.T) {
	password := []byte(strings.Repeat("a", 72))
	compressed := compressPassword(password)

	if len(compressed) != 72 {
		t.Error("compressPassword should not modify passwords exactly 72 bytes")
	}

	if string(compressed) != string(password) {
		t.Error("compressPassword modified 72-byte password content")
	}
}

func TestCompressPassword_LongPassword(t *testing.T) {
	password := []byte(strings.Repeat("a", 100))
	compressed := compressPassword(password)

	// SHA-512 produces 64 bytes
	if len(compressed) != 64 {
		t.Errorf("compressPassword should produce 64-byte hash, got %d bytes", len(compressed))
	}

	// Ensure it's different from original
	if string(compressed) == string(password) {
		t.Error("compressPassword did not hash long password")
	}
}

func TestCompressPassword_DifferentLongPasswords(t *testing.T) {
	password1 := []byte(strings.Repeat("a", 100))
	password2 := []byte(strings.Repeat("b", 100))

	compressed1 := compressPassword(password1)
	compressed2 := compressPassword(password2)

	if string(compressed1) == string(compressed2) {
		t.Error("compressPassword produced same hash for different long passwords")
	}
}

func TestHashPassword_UniqueHashes(t *testing.T) {
	password := "testpassword"

	// Generate two hashes of the same password
	hash1, err1 := HashPassword(password)
	hash2, err2 := HashPassword(password)

	if err1 != nil || err2 != nil {
		t.Fatalf("HashPassword failed: %v, %v", err1, err2)
	}

	// bcrypt includes a salt, so hashes should be different
	if hash1 == hash2 {
		t.Error("HashPassword produced identical hashes (should use salt)")
	}

	// But both should validate against the original password
	if err := CompareHash(hash1, password); err != nil {
		t.Error("First hash does not validate against password")
	}

	if err := CompareHash(hash2, password); err != nil {
		t.Error("Second hash does not validate against password")
	}
}

func TestHashPassword_EmptyPassword(t *testing.T) {
	password := ""
	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("HashPassword failed for empty password: %v", err)
	}

	err = CompareHash(hash, password)
	if err != nil {
		t.Error("CompareHash failed to match empty password")
	}
}

func TestHashPassword_SpecialCharacters(t *testing.T) {
	password := "p@ssw0rd!#$%^&*()"
	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("HashPassword failed for password with special characters: %v", err)
	}

	err = CompareHash(hash, password)
	if err != nil {
		t.Error("CompareHash failed to match password with special characters")
	}
}

func TestHashPassword_UnicodeCharacters(t *testing.T) {
	password := "пароль密码🔒"
	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("HashPassword failed for password with unicode characters: %v", err)
	}

	err = CompareHash(hash, password)
	if err != nil {
		t.Error("CompareHash failed to match password with unicode characters")
	}
}

func TestHashPassword_LongUnicodePassword(t *testing.T) {
	// Create a long password with unicode characters (over 72 bytes)
	password := strings.Repeat("密码🔒", 30) // Each character is multiple bytes

	if len([]byte(password)) <= 72 {
		t.Skip("Unicode password not long enough for test")
	}

	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("HashPassword failed for long unicode password: %v", err)
	}

	err = CompareHash(hash, password)
	if err != nil {
		t.Error("CompareHash failed to match long unicode password")
	}
}

func TestHashPassword_BcryptCostIsDefault(t *testing.T) {
	password := "testpassword"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	// Extract the cost from the hash
	cost, err := bcrypt.Cost([]byte(hash))
	if err != nil {
		t.Fatalf("Failed to extract cost from hash: %v", err)
	}

	// bcrypt.DefaultCost is 10
	if cost != bcrypt.DefaultCost {
		t.Errorf("Expected cost %d, got %d", bcrypt.DefaultCost, cost)
	}
}
