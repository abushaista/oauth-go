package http

import (
	"encoding/json"
	"net/http"

	"github.com/abushaista/oauth-go/internal/infrastructure/security"
)

// JWKSHandler handles JSON Web Key Set requests
type JWKSHandler struct {
	jwksProvider *security.JWKSProvider
}

// NewJWKSHandler creates a new JWKS HTTP handler
func NewJWKSHandler(jwksProvider *security.JWKSProvider) *JWKSHandler {
	return &JWKSHandler{
		jwksProvider: jwksProvider,
	}
}

// ServeHTTP handles GET /.well-known/jwks.json requests
func (h *JWKSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	jwks := h.jwksProvider.GetJWKS()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=300")
	json.NewEncoder(w).Encode(jwks)
}
