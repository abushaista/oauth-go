package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

// PasswordHasher handles password hashing and verification
type PasswordHasher struct {
}

// NewPasswordHasher creates a new password hasher
func NewPasswordHasher() *PasswordHasher {
	return &PasswordHasher{}
}

// preHash converts the input password into a fixed-length SHA-256 hex string.
// This avoids bcrypt's 72-byte password length limitation while preserving
// deterministic mapping for verification.
func preHash(password string) []byte {
	h := sha256.Sum256([]byte(password))
	return []byte(hex.EncodeToString(h[:]))
}

// Hash hashes a password (pre-hashes with SHA-256, then bcrypt)
func (ph *PasswordHasher) Hash(password string) (string, error) {
	toHash := preHash(password)
	bytes, err := bcrypt.GenerateFromPassword(toHash, 14)
	return string(bytes), err
}

// Verify verifies a password against a bcrypt hash (applies same pre-hash)
func (ph *PasswordHasher) Verify(hash, password string) bool {
	toVerify := preHash(password)
	err := bcrypt.CompareHashAndPassword([]byte(hash), toVerify)
	return err == nil
}

// GenerateRandomToken generates a cryptographically secure random token
func GenerateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
