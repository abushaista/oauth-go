package security

import (
	"crypto/sha256"
	"encoding/base64"
)

// PKCEValidator validates PKCE code challenges
type PKCEValidator struct{}

// NewPKCEValidator creates a new PKCE validator
func NewPKCEValidator() *PKCEValidator {
	return &PKCEValidator{}
}

// ValidateS256 validates a code verifier against a S256 code challenge
func (pv *PKCEValidator) ValidateS256(codeVerifier, codeChallenge string) bool {
	hash := sha256.Sum256([]byte(codeVerifier))
	computed := base64.RawURLEncoding.EncodeToString(hash[:])
	return computed == codeChallenge
}

// ValidatePlain validates a code verifier against a plain code challenge
func (pv *PKCEValidator) ValidatePlain(codeVerifier, codeChallenge string) bool {
	return codeVerifier == codeChallenge
}

// GenerateCodeChallenge generates an S256 code challenge
func (pv *PKCEValidator) GenerateCodeChallenge(codeVerifier string) string {
	hash := sha256.Sum256([]byte(codeVerifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}
