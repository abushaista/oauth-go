package query

import (
	"context"
	"fmt"
)

// JWKSHandler retrieves public keys in JWKS format
type JWKSHandler struct {
	// TODO: Inject JWT service to retrieve keys
}

// NewJWKSHandler creates a new JWKS handler
func NewJWKSHandler() *JWKSHandler {
	return &JWKSHandler{}
}

// Handle processes a JWKS query
func (h *JWKSHandler) Handle(ctx context.Context, q Query) (interface{}, error) {
	jwksQuery, ok := q.(*JWKSQuery)
	if !ok {
		return nil, fmt.Errorf("invalid query type")
	}

	// TODO: Retrieve keys from JWT service
	_ = jwksQuery

	return map[string]interface{}{
		"keys": []interface{}{},
	}, nil
}
