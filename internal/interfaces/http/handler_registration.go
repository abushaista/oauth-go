package http

import (
	"encoding/json"
	"net/http"
	"github.com/abushaista/oauth-go/internal/application/command"
)

// RegistrationHandler handles dynamic client registration requests
type RegistrationHandler struct {
	handler *command.RegisterClientHandler
}

// NewRegistrationHandler creates a new registration HTTP handler
func NewRegistrationHandler(handler *command.RegisterClientHandler) *RegistrationHandler {
	return &RegistrationHandler{
		handler: handler,
	}
}

// ServeHTTP handles POST /register requests
func (h *RegistrationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req command.RegisterClientCommand
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	result, err := h.handler.Handle(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}
