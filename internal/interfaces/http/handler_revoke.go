package http

import (
	"net/http"
	"github.com/abushaista/oauth-go/internal/application/command"
)

// RevokeHandler handles OAuth token revocation requests
type RevokeHandler struct {
	handler *command.RevokeHandler
}

// NewRevokeHandler creates a new revoke HTTP handler
func NewRevokeHandler(handler *command.RevokeHandler) *RevokeHandler {
	return &RevokeHandler{
		handler: handler,
	}
}

// ServeHTTP handles POST /oauth/revoke requests
func (h *RevokeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// Read client authentication
	clientID := r.PostFormValue("client_id")
	clientSecret := r.PostFormValue("client_secret")
	
	if clientID == "" {
		var ok bool
		clientID, clientSecret, ok = r.BasicAuth()
		if !ok {
			http.Error(w, "unauthorized client", http.StatusUnauthorized)
			return
		}
	}

	cmd := &command.RevokeCommand{
		Token:        r.PostFormValue("token"),
		TokenType:    r.PostFormValue("token_type_hint"),
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}

	if _, err := h.handler.Handle(r.Context(), cmd); err != nil {
		// RFC 7009: Error responses are usually not necessary if the token is just missing
		// But authentication failure should be reported.
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// RFC 7009: Successful revocation returns 200 OK
	w.WriteHeader(http.StatusOK)
}
