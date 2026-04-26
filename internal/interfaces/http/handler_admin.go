package http

import (
	"encoding/json"
	"net/http"
	"github.com/abushaista/oauth-go/internal/application/query"
)

// AdminHandler handles administrative data requests
type AdminHandler struct {
	clientQuery query.QueryHandler
	auditQuery  query.QueryHandler
}

// NewAdminHandler creates a new admin HTTP handler
func NewAdminHandler(clientQuery query.QueryHandler, auditQuery query.QueryHandler) *AdminHandler {
	return &AdminHandler{
		clientQuery: clientQuery,
		auditQuery:  auditQuery,
	}
}

// ServeHTTP handles GET /admin/api requests
func (h *AdminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	resource := r.URL.Query().Get("resource")
	var result interface{}
	var err error

	switch resource {
	case "clients":
		result, err = h.clientQuery.Handle(r.Context(), &query.ClientListQuery{})
	case "audits":
		result, err = h.auditQuery.Handle(r.Context(), &query.AuditUserQuery{UserID: "", Limit: 100})
	default:
		http.Error(w, "invalid resource", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
