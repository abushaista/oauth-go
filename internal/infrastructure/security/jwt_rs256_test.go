package security_test

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/abushaista/oauth-go/internal/infrastructure/security"
)

func newTestSigner(t *testing.T) *security.JWTSigner {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}
	return security.NewJWTSignerWithKey(key)
}

func TestJWTSigner_SignAndVerify(t *testing.T) {
	signer := newTestSigner(t)

	claims := map[string]interface{}{
		"sub":     "user-123",
		"user_id": "user-123",
		"exp":     float64(time.Now().Add(1 * time.Hour).Unix()),
	}

	token, err := signer.Sign(claims)
	if err != nil {
		t.Fatalf("Sign() returned error: %v", err)
	}

	if token == "" {
		t.Fatal("Sign() returned empty token")
	}

	// Verify the token
	verified, err := signer.Verify(token)
	if err != nil {
		t.Fatalf("Verify() returned error: %v", err)
	}

	if verified["sub"] != "user-123" {
		t.Errorf("expected sub=user-123, got %v", verified["sub"])
	}

	if verified["user_id"] != "user-123" {
		t.Errorf("expected user_id=user-123, got %v", verified["user_id"])
	}
}

func TestJWTSigner_ExpiredToken(t *testing.T) {
	signer := newTestSigner(t)

	claims := map[string]interface{}{
		"sub": "user-123",
		"exp": float64(time.Now().Add(-1 * time.Hour).Unix()), // Expired 1 hour ago
	}

	token, err := signer.Sign(claims)
	if err != nil {
		t.Fatalf("Sign() returned error: %v", err)
	}

	_, err = signer.Verify(token)
	if err == nil {
		t.Error("Verify() should return error for expired token")
	}
}

func TestJWTSigner_TamperedToken(t *testing.T) {
	signer := newTestSigner(t)

	claims := map[string]interface{}{
		"sub": "user-123",
		"exp": float64(time.Now().Add(1 * time.Hour).Unix()),
	}

	token, err := signer.Sign(claims)
	if err != nil {
		t.Fatalf("Sign() returned error: %v", err)
	}

	// Tamper with a character in the payload
	tampered := token[:len(token)/2] + "X" + token[len(token)/2+1:]

	_, err = signer.Verify(tampered)
	if err == nil {
		t.Error("Verify() should return error for tampered token")
	}
}

func TestJWTSigner_InvalidFormat(t *testing.T) {
	signer := newTestSigner(t)

	testCases := []struct {
		name  string
		token string
	}{
		{"empty string", ""},
		{"one part", "aaa"},
		{"two parts", "aaa.bbb"},
		{"four parts", "aaa.bbb.ccc.ddd"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := signer.Verify(tc.token)
			if err == nil {
				t.Error("Verify() should return error for invalid token format")
			}
		})
	}
}

func TestJWTSigner_DifferentKeysCannotVerify(t *testing.T) {
	signer1 := newTestSigner(t)
	signer2 := newTestSigner(t) // Different key pair

	claims := map[string]interface{}{
		"sub": "user-123",
		"exp": float64(time.Now().Add(1 * time.Hour).Unix()),
	}

	token, err := signer1.Sign(claims)
	if err != nil {
		t.Fatalf("Sign() returned error: %v", err)
	}

	_, err = signer2.Verify(token)
	if err == nil {
		t.Error("Verify() should fail when using a different key pair")
	}
}

func TestJWTSigner_AutoSetsIssuedAt(t *testing.T) {
	signer := newTestSigner(t)

	before := time.Now().Unix()

	claims := map[string]interface{}{
		"sub": "user-123",
		"exp": float64(time.Now().Add(1 * time.Hour).Unix()),
	}

	token, err := signer.Sign(claims)
	if err != nil {
		t.Fatalf("Sign() returned error: %v", err)
	}

	after := time.Now().Unix()

	verified, err := signer.Verify(token)
	if err != nil {
		t.Fatalf("Verify() returned error: %v", err)
	}

	iat, ok := verified["iat"].(float64)
	if !ok {
		t.Fatal("iat claim should be present and numeric")
	}

	if int64(iat) < before || int64(iat) > after {
		t.Errorf("iat should be between %d and %d, got %d", before, after, int64(iat))
	}
}

func TestJWTSigner_WithJWKSProvider(t *testing.T) {
	provider := security.NewJWKSProvider()
	signer := security.NewJWTSigner(provider)

	claims := map[string]interface{}{
		"sub": "user-456",
		"exp": float64(time.Now().Add(1 * time.Hour).Unix()),
	}

	token, err := signer.Sign(claims)
	if err != nil {
		t.Fatalf("Sign() returned error: %v", err)
	}

	verified, err := signer.Verify(token)
	if err != nil {
		t.Fatalf("Verify() returned error: %v", err)
	}

	if verified["sub"] != "user-456" {
		t.Errorf("expected sub=user-456, got %v", verified["sub"])
	}
}
