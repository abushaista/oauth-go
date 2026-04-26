package http

import (
	"encoding/json"
	"net/http"

	"github.com/abushaista/oauth-go/internal/application/command"
)

// TokenHandler handles OAuth token requests
type TokenHandler struct {
	tokenHandler   *command.TokenHandler
	refreshHandler *command.RefreshHandler
}

// NewTokenHandler creates a new token HTTP handler
func NewTokenHandler(tokenHandler *command.TokenHandler, refreshHandler *command.RefreshHandler) *TokenHandler {
	return &TokenHandler{
		tokenHandler:   tokenHandler,
		refreshHandler: refreshHandler,
	}
}

// ServeHTTP handles POST /oauth/token requests
func (h *TokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	grantType := r.PostFormValue("grant_type")
	var result interface{}
	var err error

	if grantType == "authorization_code" {
		// Read client authentication
		clientID := r.PostFormValue("client_id")
		clientSecret := r.PostFormValue("client_secret")
		
		// Fallback to basic auth as per RFC 6749
		if clientID == "" {
			var ok bool
			clientID, clientSecret, ok = r.BasicAuth()
			if !ok {
				http.Error(w, "unauthorized client", http.StatusUnauthorized)
				return
			}
		}

		cmd := &command.TokenCommand{
			GrantType:    grantType,
			Code:         r.PostFormValue("code"),
			CodeVerifier: r.PostFormValue("code_verifier"),
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURI:  r.PostFormValue("redirect_uri"),
			Scope:        r.PostFormValue("scope"),
		}
		result, err = h.tokenHandler.Handle(r.Context(), cmd)
	} else if grantType == "refresh_token" {
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

		cmd := &command.RefreshCommand{
			GrantType:    grantType,
			RefreshToken: r.PostFormValue("refresh_token"),
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Scope:        r.PostFormValue("scope"),
		}
		result, err = h.refreshHandler.Handle(r.Context(), cmd)
	} else {
		http.Error(w, "unsupported grant type", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")

	json.NewEncoder(w).Encode(result)
}
