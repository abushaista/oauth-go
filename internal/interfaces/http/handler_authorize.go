package http

import (
	"net/http"
	"net/url"

	"github.com/abushaista/oauth-go/internal/application/command"
)

// AuthorizeHandler handles OAuth authorization requests
type AuthorizeHandler struct {
	handler        *command.AuthorizeHandler
	sessionManager *SessionManager
}

// NewAuthorizeHandler creates a new authorize HTTP handler
func NewAuthorizeHandler(handler *command.AuthorizeHandler, sessionManager *SessionManager) *AuthorizeHandler {
	return &AuthorizeHandler{
		handler:        handler,
		sessionManager: sessionManager,
	}
}

// ServeHTTP handles GET /oauth/authorize requests
func (h *AuthorizeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	
	params := r.Form
	if r.Method == http.MethodPost {
		params = r.PostForm
	}

	sessionCookie, err := r.Cookie("session")
	if err != nil || sessionCookie.Value == "" {
		returnTo := url.QueryEscape(r.URL.RequestURI())
		http.Redirect(w, r, "/ui/login.html?return_to="+returnTo, http.StatusFound)
		return
	}

	userID, err := h.sessionManager.VerifySession(sessionCookie.Value)
	if err != nil {
		returnTo := url.QueryEscape(r.URL.RequestURI())
		http.Redirect(w, r, "/ui/login.html?return_to="+returnTo, http.StatusFound)
		return
	}

	cmd := &command.AuthorizeCommand{
		ClientID:            params.Get("client_id"),
		RedirectURI:         params.Get("redirect_uri"),
		ResponseType:        params.Get("response_type"),
		Scope:               params.Get("scope"),
		State:               params.Get("state"),
		CodeChallenge:       params.Get("code_challenge"),
		CodeChallengeMethod: params.Get("code_challenge_method"),
		UserID:              userID,
	}

	result, err := h.handler.Handle(r.Context(), cmd)
	if err != nil {
		redirectURI := params.Get("redirect_uri")
		if redirectURI == "" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		u, _ := url.Parse(redirectURI)
		q := u.Query()
		q.Set("error", "invalid_request")
		q.Set("error_description", err.Error())
		if state := params.Get("state"); state != "" {
			q.Set("state", state)
		}
		u.RawQuery = q.Encode()
		http.Redirect(w, r, u.String(), http.StatusFound)
		return
	}

	resMap := result.(map[string]interface{})
	if reqConsent, ok := resMap["require_consent"].(bool); ok && reqConsent {
		// Ensure we preserve all query parameters for the consent screen
		consentURL := "/consent?" + r.URL.RawQuery
		http.Redirect(w, r, consentURL, http.StatusFound)
		return
	}

	code := resMap["code"].(string)
	state := resMap["state"].(string)

	redirectURI := params.Get("redirect_uri")
	u, _ := url.Parse(redirectURI)
	q := u.Query()
	q.Set("code", code)
	if state != "" {
		q.Set("state", state)
	}
	u.RawQuery = q.Encode()

	http.Redirect(w, r, u.String(), http.StatusFound)
}
