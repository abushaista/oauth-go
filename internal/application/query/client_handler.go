package query

import (
	"context"
	"github.com/abushaista/oauth-go/internal/domain"
)

// ClientListQuery handles requests to list all clients
type ClientListQuery struct{}

// ClientQueryHandler retrieves client information
type ClientQueryHandler struct {
	clientRepo domain.ClientRepository
}

// NewClientQueryHandler creates a new client query handler
func NewClientQueryHandler(clientRepo domain.ClientRepository) *ClientQueryHandler {
	return &ClientQueryHandler{
		clientRepo: clientRepo,
	}
}

// Handle processes client-related queries
func (h *ClientQueryHandler) Handle(ctx context.Context, q Query) (interface{}, error) {
	switch q.(type) {
	case *ClientListQuery:
		return h.clientRepo.FindAll(ctx)
	default:
		return nil, nil
	}
}
