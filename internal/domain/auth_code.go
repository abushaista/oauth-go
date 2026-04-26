package domain

import (
	"context"
	"time"
)

type AuthorizationCode struct {
	Code                string
	UserID              string
	ClientID            string
	CodeChallenge       string
	CodeChallengeMethod string
	Scope               string
	ExpiresAt           time.Time
}

// AuthorizationCodeRepository defines the interface for authorization code operations
type AuthorizationCodeRepository interface {
	// Create stores a new authorization code
	Create(ctx context.Context, authCode *AuthorizationCode) error

	// FindByCode retrieves an authorization code by its code value
	FindByCode(ctx context.Context, code string) (*AuthorizationCode, error)

	// Delete removes an authorization code
	Delete(ctx context.Context, code string) error
}
