package http

import (
	"net/http"

	"github.com/abushaista/oauth-go/internal/domain"
)

// RoleMiddleware enforces Role-Based Access Control
type RoleMiddleware struct {
	userRepo domain.UserRepository
}

// NewRoleMiddleware creates a new RoleMiddleware
func NewRoleMiddleware(userRepo domain.UserRepository) *RoleMiddleware {
	return &RoleMiddleware{
		userRepo: userRepo,
	}
}

// RequireRole ensures the user has one of the specified roles
func (m *RoleMiddleware) RequireRole(requiredRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract user_id from context (injected by AuthMiddleware)
			userID, ok := r.Context().Value(UserIDKey).(string)
			if !ok || userID == "" {
				http.Error(w, "unauthorized: missing user information", http.StatusUnauthorized)
				return
			}

			// Fetch user from DB to get current role
			user, err := m.userRepo.FindByID(r.Context(), userID)
			if err != nil || user == nil {
				http.Error(w, "unauthorized: user not found", http.StatusUnauthorized)
				return
			}

			// Check if user has required role
			hasRole := false
			for _, role := range requiredRoles {
				if user.Role == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				http.Error(w, "forbidden: insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
