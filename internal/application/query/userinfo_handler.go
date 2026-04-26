package query

import (
	"context"
	"fmt"
	"github.com/abushaista/oauth-go/internal/domain"
)

// UserinfoQuery handles OIDC userinfo requests
type UserinfoQuery struct {
	UserID string
}

// UserinfoHandler retrieves user profile information
type UserinfoHandler struct {
	userRepo domain.UserRepository
}

// NewUserinfoHandler creates a new userinfo handler
func NewUserinfoHandler(userRepo domain.UserRepository) *UserinfoHandler {
	return &UserinfoHandler{
		userRepo: userRepo,
	}
}

// Handle processes a userinfo query
func (h *UserinfoHandler) Handle(ctx context.Context, q Query) (interface{}, error) {
	query, ok := q.(*UserinfoQuery)
	if !ok {
		return nil, fmt.Errorf("invalid query type")
	}

	user, err := h.userRepo.FindByID(ctx, query.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// OIDC standard claims
	return map[string]interface{}{
		"sub":      user.ID,
		"username": user.Username,
		"name":     user.Username, // Placeholder
		"email":    user.Username + "@example.com", // Placeholder
	}, nil
}
