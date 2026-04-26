package http

import (
	"context"
	"net/http"
	"strings"
	"github.com/abushaista/oauth-go/internal/domain"
)

type contextKey string
const UserIDKey contextKey = "user_id"

// AuthMiddleware validates the Bearer access token
type AuthMiddleware struct {
	tokenRepo domain.AccessTokenRepository
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(tokenRepo domain.AccessTokenRepository) *AuthMiddleware {
	return &AuthMiddleware{
		tokenRepo: tokenRepo,
	}
}

// Wrap executes the authentication logic
func (m *AuthMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		token := parts[1]
		accessToken, err := m.tokenRepo.FindByToken(r.Context(), token)
		if err != nil || accessToken == nil {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Inject user_id into context
		ctx := context.WithValue(r.Context(), UserIDKey, accessToken.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
