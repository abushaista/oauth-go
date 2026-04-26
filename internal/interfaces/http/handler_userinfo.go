package http

import (
	"encoding/json"
	"net/http"
	"github.com/abushaista/oauth-go/internal/application/query"
)

// UserinfoHandler handles Userinfo requests
type UserinfoHandler struct {
	queryHandler *query.UserinfoHandler
}

// NewUserinfoHandler creates a new userinfo HTTP handler
func NewUserinfoHandler(queryHandler *query.UserinfoHandler) *UserinfoHandler {
	return &UserinfoHandler{
		queryHandler: queryHandler,
	}
}

// ServeHTTP handles GET /userinfo requests
func (h *UserinfoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value(UserIDKey).(string)
	if !ok || userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	q := &query.UserinfoQuery{UserID: userID}
	result, err := h.queryHandler.Handle(r.Context(), q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
