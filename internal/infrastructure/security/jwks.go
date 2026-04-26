package security

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"math/big"
)

// JWKS represents the JSON Web Key Set
type JWKS struct {
	Keys []JWK `json:"keys"`
}

// JWK represents a single JSON Web Key
type JWK struct {
	Kty string `json:"kty"`           // Key type (RSA, EC, etc.)
	Kid string `json:"kid"`           // Key ID
	Use string `json:"use"`           // Key usage (sig, enc)
	Alg string `json:"alg"`           // Algorithm
	N   string `json:"n,omitempty"`   // RSA modulus
	E   string `json:"e,omitempty"`   // RSA exponent
	Crv string `json:"crv,omitempty"` // EC curve
	X   string `json:"x,omitempty"`   // EC X coordinate
	Y   string `json:"y,omitempty"`   // EC Y coordinate
}

// JWKSProvider provides the JSON Web Key Set
type JWKSProvider struct {
	PrivateKey *rsa.PrivateKey
}

// NewJWKSProvider creates a new JWKS provider
func NewJWKSProvider() *JWKSProvider {
	// Generate a dynamic RSA key-pair on startup 
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic("failed to generate RSA key for JWKS: " + err.Error())
	}

	return &JWKSProvider{
		PrivateKey: priv,
	}
}

// GetJWKS returns the JSON Web Key Set
func (jp *JWKSProvider) GetJWKS() *JWKS {
	pub := jp.PrivateKey.PublicKey

	nBytes := pub.N.Bytes()
	eBytes := big.NewInt(int64(pub.E)).Bytes()

	return &JWKS{
		Keys: []JWK{
			{
				Kty: "RSA",
				Kid: "global-rs256-key-1",
				Use: "sig",
				Alg: "RS256",
				N:   base64.RawURLEncoding.EncodeToString(nBytes),
				E:   base64.RawURLEncoding.EncodeToString(eBytes),
			},
		},
	}
}
