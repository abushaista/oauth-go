package security_test

import (
	"testing"

	"github.com/abushaista/oauth-go/internal/infrastructure/security"
)

func TestPasswordHasher_HashAndVerify(t *testing.T) {
	hasher := security.NewPasswordHasher()

	password := "my-secure-password-123"

	hash, err := hasher.Hash(password)
	if err != nil {
		t.Fatalf("Hash() returned error: %v", err)
	}

	if hash == "" {
		t.Fatal("Hash() returned empty string")
	}

	if hash == password {
		t.Fatal("Hash() returned the plaintext password")
	}

	// Verify correct password
	if !hasher.Verify(hash, password) {
		t.Error("Verify() returned false for correct password")
	}

	// Verify wrong password
	if hasher.Verify(hash, "wrong-password") {
		t.Error("Verify() returned true for wrong password")
	}
}

func TestPasswordHasher_DifferentHashesForSamePassword(t *testing.T) {
	hasher := security.NewPasswordHasher()

	hash1, _ := hasher.Hash("same-password")
	hash2, _ := hasher.Hash("same-password")

	if hash1 == hash2 {
		t.Error("Two calls to Hash() with the same password should produce different hashes (bcrypt salt)")
	}
}

func TestPasswordHasher_EmptyPassword(t *testing.T) {
	hasher := security.NewPasswordHasher()

	hash, err := hasher.Hash("")
	if err != nil {
		t.Fatalf("Hash() should accept empty password: %v", err)
	}

	if !hasher.Verify(hash, "") {
		t.Error("Verify() should accept empty password that was hashed")
	}
}

func TestPasswordHasher_LongPassword(t *testing.T) {
	hasher := security.NewPasswordHasher()

	// bcrypt has a 72-byte limit; this tests behavior near that boundary
	longPassword := "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*()_+"

	hash, err := hasher.Hash(longPassword)
	if err != nil {
		t.Fatalf("Hash() with long password returned error: %v", err)
	}

	if !hasher.Verify(hash, longPassword) {
		t.Error("Verify() failed for long password")
	}
}

func TestGenerateRandomToken(t *testing.T) {
	token1, err := security.GenerateRandomToken(32)
	if err != nil {
		t.Fatalf("GenerateRandomToken() returned error: %v", err)
	}

	token2, err := security.GenerateRandomToken(32)
	if err != nil {
		t.Fatalf("GenerateRandomToken() returned error: %v", err)
	}

	if token1 == "" || token2 == "" {
		t.Error("GenerateRandomToken() returned empty string")
	}

	if token1 == token2 {
		t.Error("Two calls to GenerateRandomToken() should produce different tokens")
	}
}

func TestGenerateRandomToken_DifferentLengths(t *testing.T) {
	short, _ := security.GenerateRandomToken(8)
	long, _ := security.GenerateRandomToken(64)

	if len(short) >= len(long) {
		t.Error("Longer byte input should produce longer base64 output")
	}
}
