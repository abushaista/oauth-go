package query

import (
	"context"
)

// Query is a marker interface for all queries
type Query interface{}

// QueryHandler is the interface all query handlers must implement
type QueryHandler interface {
	Handle(ctx context.Context, q Query) (interface{}, error)
}

// JWKSQuery retrieves the JSON Web Key Set
type JWKSQuery struct {
	KeyID string // Optional: specific key ID
}

// AuditUserQuery fetches audit logs for a specific user ID
type AuditUserQuery struct {
	UserID string
	Limit  int
}

// AuditClientQuery fetches audit logs for a specific client ID
type AuditClientQuery struct {
	ClientID string
	Limit    int
}
