package command

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/abushaista/oauth-go/internal/infrastructure/security"
)

// generateRandomToken creates a secure random token string using hex encoding
func generateRandomToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		// Fallback if crypto/rand fails
		return hex.EncodeToString([]byte(time.Now().String()))
	}
	return hex.EncodeToString(b)
}

// generateRandomCode creates a secure random authorization code
func generateRandomCode() string {
	token, _ := security.GenerateRandomToken(24)
	return token
}

// getExpirationTime returns a time object 1 hour in the future
func getExpirationTime() time.Time {
	return time.Now().Add(1 * time.Hour)
}
