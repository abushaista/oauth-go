package http

import (
	"encoding/json"
	"net/http"

	"github.com/abushaista/oauth-go/internal/application/command"
)

// LoginHandler handles login requests
type LoginHandler struct {
	handler        *command.LoginHandler
	sessionManager *SessionManager
}

// NewLoginHandler creates a new login HTTP handler
func NewLoginHandler(handler *command.LoginHandler, sessionManager *SessionManager) *LoginHandler {
	return &LoginHandler{
		handler:        handler,
		sessionManager: sessionManager,
	}
}

// ServeHTTP handles POST /login requests
func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var username, password string

	if r.Header.Get("Content-Type") == "application/json" {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err == nil {
			username = req.Username
			password = req.Password
		}
	} else {
		r.ParseForm()
		username = r.FormValue("username")
		password = r.FormValue("password")
	}

	if username == "" || password == "" {
		http.Error(w, "Missing credentials", http.StatusBadRequest)
		return
	}

	cmd := &command.LoginCommand{
		Username: username,
		Password: password,
	}

	result, err := h.handler.Handle(r.Context(), cmd)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	resMap, ok := result.(map[string]interface{})
	if !ok || !resMap["success"].(bool) {
		http.Error(w, "Login failed", http.StatusUnauthorized)
		return
	}

	userID := resMap["user_id"].(string)

	// Use SessionManager to create session and set cookie
	token, err := h.sessionManager.CreateSession(userID)
	if err != nil {
		http.Error(w, "failed to create session", http.StatusInternalServerError)
		return
	}

	h.sessionManager.SetSessionCookie(w, token)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Login successful",
		"user_id": userID,
	})
}
