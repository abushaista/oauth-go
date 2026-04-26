package http

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/abushaista/oauth-go/internal/domain"
)

// AuditMiddleware logs generalized HTTP requests for telemetry
type AuditMiddleware struct {
	auditRepo domain.AuditRepository
}

// NewAuditMiddleware creates a new audit middleware
func NewAuditMiddleware(auditRepo domain.AuditRepository) *AuditMiddleware {
	return &AuditMiddleware{
		auditRepo: auditRepo,
	}
}

func generateLogID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// Wrap executes the middleware logic around the next handler
func (m *AuditMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Execute the wrapped handler
		next.ServeHTTP(w, r)
		
		// Collect best-effort telemetry
		clientID := r.URL.Query().Get("client_id")
		if clientID == "" {
			clientID = r.PostFormValue("client_id")
		}

		if clientID == "" {
			clientID = "system"
		}

		auditEntry := &domain.Audit{
			ID:        generateLogID(),
			UserID:    "anonymous",
			ClientID:  clientID,
			Action:    "HTTP_REQUEST",
			Details:   fmt.Sprintf("Method: %s | Path: %s", r.Method, r.URL.Path),
			IPAddress: r.RemoteAddr,
			CreatedAt: time.Now(),
		}

		_ = m.auditRepo.Create(r.Context(), auditEntry)
	})
}
