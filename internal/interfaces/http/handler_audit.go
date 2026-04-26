package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/abushaista/oauth-go/internal/application/query"
)

// AuditHandler surfaces audit event telemetry
type AuditHandler struct {
	queryHandler *query.AuditQueryHandler
}

// NewAuditHandler creates an HTTP endpoint handler for Audits
func NewAuditHandler(queryHandler *query.AuditQueryHandler) *AuditHandler {
	return &AuditHandler{
		queryHandler: queryHandler,
	}
}

// ServeHTTP handles GET /audits requests
func (h *AuditHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("user_id")
	clientID := r.URL.Query().Get("client_id")
	limitStr := r.URL.Query().Get("limit")

	limit := 50
	if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 {
		limit = parsed
	}

	var cqrsQuery query.Query

	if userID != "" {
		cqrsQuery = &query.AuditUserQuery{
			UserID: userID,
			Limit:  limit,
		}
	} else if clientID != "" {
		cqrsQuery = &query.AuditClientQuery{
			ClientID: clientID,
			Limit:    limit,
		}
	} else {
		http.Error(w, "must provide user_id or client_id", http.StatusBadRequest)
		return
	}

	result, err := h.queryHandler.Handle(r.Context(), cqrsQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
