package security

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// JWTSigner handles JWT signing and verification with RS256
type JWTSigner struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	keyID      string
}

// NewJWTSigner creates a new JWT signer using a shared key from JWKSProvider
func NewJWTSigner(provider *JWKSProvider) *JWTSigner {
	return &JWTSigner{
		privateKey: provider.PrivateKey,
		publicKey:  &provider.PrivateKey.PublicKey,
		keyID:      "global-rs256-key-1",
	}
}

// NewJWTSignerWithKey creates a signer with a specific RSA key pair (useful for testing)
func NewJWTSignerWithKey(privateKey *rsa.PrivateKey) *JWTSigner {
	return &JWTSigner{
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
		keyID:      "global-rs256-key-1",
	}
}

// Sign creates a signed JWT token with RS256
func (js *JWTSigner) Sign(claims map[string]interface{}) (string, error) {
	// Build JWT header
	header := map[string]string{
		"alg": "RS256",
		"typ": "JWT",
		"kid": js.keyID,
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("failed to marshal header: %w", err)
	}

	// Set issued-at if not present
	if _, ok := claims["iat"]; !ok {
		claims["iat"] = time.Now().Unix()
	}

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("failed to marshal claims: %w", err)
	}

	// Encode header and payload
	headerB64 := base64.RawURLEncoding.EncodeToString(headerJSON)
	payloadB64 := base64.RawURLEncoding.EncodeToString(claimsJSON)

	// Create signing input
	signingInput := headerB64 + "." + payloadB64

	// Sign with RS256 (RSASSA-PKCS1-v1_5 with SHA-256)
	hash := sha256.Sum256([]byte(signingInput))
	signature, err := rsa.SignPKCS1v15(rand.Reader, js.privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	signatureB64 := base64.RawURLEncoding.EncodeToString(signature)

	return signingInput + "." + signatureB64, nil
}

// Verify validates a JWT token and returns its claims
func (js *JWTSigner) Verify(token string) (map[string]interface{}, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format: expected 3 parts, got %d", len(parts))
	}

	headerB64 := parts[0]
	payloadB64 := parts[1]
	signatureB64 := parts[2]

	// Decode and validate header
	headerJSON, err := base64.RawURLEncoding.DecodeString(headerB64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode header: %w", err)
	}

	var header map[string]string
	if err := json.Unmarshal(headerJSON, &header); err != nil {
		return nil, fmt.Errorf("failed to parse header: %w", err)
	}

	if header["alg"] != "RS256" {
		return nil, fmt.Errorf("unsupported algorithm: %s", header["alg"])
	}

	// Verify signature
	signingInput := headerB64 + "." + payloadB64
	hash := sha256.Sum256([]byte(signingInput))

	signature, err := base64.RawURLEncoding.DecodeString(signatureB64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode signature: %w", err)
	}

	if err := rsa.VerifyPKCS1v15(js.publicKey, crypto.SHA256, hash[:], signature); err != nil {
		return nil, fmt.Errorf("invalid signature: %w", err)
	}

	// Decode claims
	claimsJSON, err := base64.RawURLEncoding.DecodeString(payloadB64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode claims: %w", err)
	}

	var claims map[string]interface{}
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		return nil, fmt.Errorf("failed to parse claims: %w", err)
	}

	// Check expiration
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return nil, fmt.Errorf("token has expired")
		}
	}

	return claims, nil
}
