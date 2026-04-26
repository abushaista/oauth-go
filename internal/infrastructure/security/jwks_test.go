package security_test

import (
	"encoding/json"
	"testing"

	"github.com/abushaista/oauth-go/internal/infrastructure/security"
)

func TestJWKSProvider_GetJWKS(t *testing.T) {
	provider := security.NewJWKSProvider()
	jwks := provider.GetJWKS()

	if jwks == nil {
		t.Fatal("GetJWKS() returned nil")
	}

	if len(jwks.Keys) != 1 {
		t.Fatalf("expected 1 key, got %d", len(jwks.Keys))
	}

	key := jwks.Keys[0]

	if key.Kty != "RSA" {
		t.Errorf("expected kty=RSA, got %s", key.Kty)
	}

	if key.Alg != "RS256" {
		t.Errorf("expected alg=RS256, got %s", key.Alg)
	}

	if key.Use != "sig" {
		t.Errorf("expected use=sig, got %s", key.Use)
	}

	if key.Kid == "" {
		t.Error("kid should not be empty")
	}

	if key.N == "" {
		t.Error("RSA modulus (N) should not be empty")
	}

	if key.E == "" {
		t.Error("RSA exponent (E) should not be empty")
	}
}

func TestJWKSProvider_JSONSerialization(t *testing.T) {
	provider := security.NewJWKSProvider()
	jwks := provider.GetJWKS()

	data, err := json.Marshal(jwks)
	if err != nil {
		t.Fatalf("JSON marshaling failed: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("JSON unmarshaling failed: %v", err)
	}

	keys, ok := parsed["keys"].([]interface{})
	if !ok || len(keys) != 1 {
		t.Fatal("expected 'keys' array with 1 element")
	}

	keyMap := keys[0].(map[string]interface{})

	requiredFields := []string{"kty", "kid", "use", "alg", "n", "e"}
	for _, field := range requiredFields {
		if _, ok := keyMap[field]; !ok {
			t.Errorf("missing required JWKS field: %s", field)
		}
	}
}

func TestJWKSProvider_ConsistentKeys(t *testing.T) {
	provider := security.NewJWKSProvider()

	jwks1 := provider.GetJWKS()
	jwks2 := provider.GetJWKS()

	if jwks1.Keys[0].N != jwks2.Keys[0].N {
		t.Error("multiple calls to GetJWKS() should return the same key material")
	}
}
