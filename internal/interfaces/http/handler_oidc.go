package http

import (
	"encoding/json"
	"net/http"
	"os"
)

// OIDCProviderConfig represents the OIDC provider metadata
type OIDCProviderConfig struct {
	Issuer                            string   `json:"issuer"`
	AuthorizationEndpoint             string   `json:"authorization_endpoint"`
	TokenEndpoint                     string   `json:"token_endpoint"`
	UserinfoEndpoint                  string   `json:"userinfo_endpoint"`
	JwksURI                           string   `json:"jwks_uri"`
	RevocationEndpoint                string   `json:"revocation_endpoint"`
	ResponseTypesSupported            []string `json:"response_types_supported"`
	SubjectTypesSupported             []string `json:"subject_types_supported"`
	IDTokenSigningAlgValuesSupported  []string `json:"id_token_signing_alg_values_supported"`
	ScopesSupported                   []string `json:"scopes_supported"`
	TokenEndpointAuthMethodsSupported []string `json:"token_endpoint_auth_methods_supported"`
	ClaimsSupported                   []string `json:"claims_supported"`
}

// OIDCHandler handles OIDC discovery and related endpoints
type OIDCHandler struct {
	issuer string
}

// NewOIDCHandler creates a new OIDC handler
func NewOIDCHandler() *OIDCHandler {
	issuer := os.Getenv("OIDC_ISSUER")
	if issuer == "" {
		issuer = "http://localhost:8080"
	}
	return &OIDCHandler{
		issuer: issuer,
	}
}

// ServeHTTP handles /.well-known/openid-configuration
func (h *OIDCHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	config := OIDCProviderConfig{
		Issuer:                            h.issuer,
		AuthorizationEndpoint:             h.issuer + "/oauth/authorize",
		TokenEndpoint:                     h.issuer + "/oauth/token",
		UserinfoEndpoint:                  h.issuer + "/userinfo",
		JwksURI:                           h.issuer + "/.well-known/jwks.json",
		RevocationEndpoint:                h.issuer + "/oauth/revoke",
		ResponseTypesSupported:            []string{"code"},
		SubjectTypesSupported:             []string{"public"},
		IDTokenSigningAlgValuesSupported:  []string{"RS256"},
		ScopesSupported:                   []string{"openid", "profile", "email", "offline_access"},
		TokenEndpointAuthMethodsSupported: []string{"client_secret_post", "client_secret_basic"},
		ClaimsSupported:                   []string{"sub", "iss", "iat", "exp", "name", "email"},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}
