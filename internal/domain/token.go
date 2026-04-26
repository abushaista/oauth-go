package domain

import (
	"context"
	"time"
)

type AccessToken struct {
	Token     string
	UserID    string
	ClientID  string
	Scope     string
	ExpiresAt time.Time
}

type RefreshToken struct {
	Token     string
	UserID    string
	ClientID  string
	ExpiresAt time.Time
	Revoked   bool
}

// AccessTokenRepository defines the interface for access token operations
type AccessTokenRepository interface {
	// Create stores a new access token
	Create(ctx context.Context, token *AccessToken) error

	// FindByToken retrieves an access token by its token value
	FindByToken(ctx context.Context, token string) (*AccessToken, error)

	// Delete removes an access token
	Delete(ctx context.Context, token string) error
}

// RefreshTokenRepository defines the interface for refresh token operations
type RefreshTokenRepository interface {
	// Create stores a new refresh token
	Create(ctx context.Context, token *RefreshToken) error

	// FindByToken retrieves a refresh token by its token value
	FindByToken(ctx context.Context, token string) (*RefreshToken, error)

	// Revoke marks a refresh token as revoked
	Revoke(ctx context.Context, token string) error

	// Delete removes a refresh token
	Delete(ctx context.Context, token string) error
}
