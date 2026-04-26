package domain

import (
	"crypto/sha256"
	"encoding/base64"
)

type PKCE struct {
	CodeChallenge       string
	CodeChallengeMethod string // "plain" or "S256"
}

// ValidatePKCE validates the code verifier against the code challenge
func (p *PKCE) ValidatePKCE(codeVerifier string) bool {
	if p.CodeChallengeMethod == "plain" {
		return codeVerifier == p.CodeChallenge
	}
	if p.CodeChallengeMethod == "S256" {
		hasher := sha256.New()
		hasher.Write([]byte(codeVerifier))
		hash := base64.RawURLEncoding.EncodeToString(hasher.Sum(nil))
		return hash == p.CodeChallenge
	}
	return false
}
