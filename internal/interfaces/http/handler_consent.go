package http

import (
	"encoding/json"
	"net/http"

	"github.com/abushaista/oauth-go/internal/application/command"
)

// ConsentHandler handles user consent for client applications
type ConsentHandler struct {
	handler        *command.ConsentHandler
	sessionManager *SessionManager
}

// NewConsentHandler creates a new consent HTTP handler
func NewConsentHandler(handler *command.ConsentHandler, sessionManager *SessionManager) *ConsentHandler {
	return &ConsentHandler{
		handler:        handler,
		sessionManager: sessionManager,
	}
}

// ServeHTTP handles GET/POST /consent requests
func (h *ConsentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.showConsentPage(w, r)
	case http.MethodPost:
		h.processConsent(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// showConsentPage redirects to the Vue consent UI with query params
func (h *ConsentHandler) showConsentPage(w http.ResponseWriter, r *http.Request) {
	// Forward all OAuth query params to the consent UI
	consentURL := "/ui/consent.html?" + r.URL.RawQuery
	http.Redirect(w, r, consentURL, http.StatusFound)
}

// processConsent stores the user's consent decision
func (h *ConsentHandler) processConsent(w http.ResponseWriter, r *http.Request) {
	// Read user ID from session cookie securely using SessionManager
	sessionCookie, err := r.Cookie("session")
	if err != nil || sessionCookie.Value == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := h.sessionManager.VerifySession(sessionCookie.Value)
	if err != nil {
		http.Error(w, "invalid session: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Parse request body
	var req struct {
		ClientID string   `json:"client_id"`
		Scopes   []string `json:"scopes"`
		Approved bool     `json:"approved"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if !req.Approved {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "consent denied",
		})
		return
	}

	// Store consent via command handler
	cmd := &command.ConsentCommand{
		UserID:   userID,
		ClientID: req.ClientID,
		Scopes:   req.Scopes,
	}

	if _, err := h.handler.Handle(r.Context(), cmd); err != nil {
		http.Error(w, "failed to store consent: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "consent granted",
	})
}
