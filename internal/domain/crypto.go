package domain

// Crypto defines cryptographic constants and utilities
const (
	// PKCE methods
	PKCEMethodPlain = "plain"
	PKCEMethodS256  = "S256"

	// Token types
	TokenTypeBearer = "Bearer"

	// OAuth scopes
	ScopeOpenID        = "openid"
	ScopeProfile       = "profile"
	ScopeEmail         = "email"
	ScopeOfflineAccess = "offline_access"
)
