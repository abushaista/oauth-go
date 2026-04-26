package domain_test

import (
	"crypto/sha256"
	"encoding/base64"
	"testing"

	"github.com/abushaista/oauth-go/internal/domain"
)

func TestPKCE_ValidatePlain(t *testing.T) {
	pkce := &domain.PKCE{
		CodeChallenge:       "my-plain-verifier",
		CodeChallengeMethod: "plain",
	}

	if !pkce.ValidatePKCE("my-plain-verifier") {
		t.Error("ValidatePKCE() should return true for matching plain verifier")
	}

	if pkce.ValidatePKCE("wrong-verifier") {
		t.Error("ValidatePKCE() should return false for mismatched plain verifier")
	}
}

func TestPKCE_ValidateS256(t *testing.T) {
	verifier := "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"

	// Generate challenge: BASE64URL(SHA256(verifier))
	hasher := sha256.New()
	hasher.Write([]byte(verifier))
	challenge := base64.RawURLEncoding.EncodeToString(hasher.Sum(nil))

	pkce := &domain.PKCE{
		CodeChallenge:       challenge,
		CodeChallengeMethod: "S256",
	}

	if !pkce.ValidatePKCE(verifier) {
		t.Error("ValidatePKCE() should return true for correct S256 verifier")
	}

	if pkce.ValidatePKCE("wrong-verifier-value") {
		t.Error("ValidatePKCE() should return false for wrong S256 verifier")
	}
}

func TestPKCE_UnsupportedMethod(t *testing.T) {
	pkce := &domain.PKCE{
		CodeChallenge:       "anything",
		CodeChallengeMethod: "unsupported",
	}

	if pkce.ValidatePKCE("anything") {
		t.Error("ValidatePKCE() should return false for unsupported method")
	}
}

func TestPKCE_EmptyMethod(t *testing.T) {
	pkce := &domain.PKCE{
		CodeChallenge:       "value",
		CodeChallengeMethod: "",
	}

	if pkce.ValidatePKCE("value") {
		t.Error("ValidatePKCE() should return false for empty method")
	}
}
