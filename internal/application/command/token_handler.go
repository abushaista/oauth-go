package command

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/abushaista/oauth-go/internal/domain"
)

type IDTokenSigner interface {
	Sign(claims map[string]interface{}) (string, error)
}

// TokenHandler exchanges authorization codes and refresh tokens for access tokens
type TokenHandler struct {
	authCodeRepo domain.AuthorizationCodeRepository
	tokenRepo    domain.AccessTokenRepository
	refreshRepo  domain.RefreshTokenRepository
	clientRepo   domain.ClientRepository
	auditRepo    domain.AuditRepository
	jwtSigner    IDTokenSigner
}

// NewTokenHandler creates a new token handler
func NewTokenHandler(
	authCodeRepo domain.AuthorizationCodeRepository,
	tokenRepo domain.AccessTokenRepository,
	refreshRepo domain.RefreshTokenRepository,
	clientRepo domain.ClientRepository,
	auditRepo domain.AuditRepository,
	jwtSigner IDTokenSigner,
) *TokenHandler {
	return &TokenHandler{
		authCodeRepo: authCodeRepo,
		tokenRepo:    tokenRepo,
		refreshRepo:  refreshRepo,
		clientRepo:   clientRepo,
		auditRepo:    auditRepo,
		jwtSigner:    jwtSigner,
	}
}

// Handle processes a token request
func (h *TokenHandler) Handle(ctx context.Context, cmd Command) (interface{}, error) {
	tokenCmd, ok := cmd.(*TokenCommand)
	if !ok {
		return nil, fmt.Errorf("invalid command type")
	}

	// Validate client credentials
	client, err := h.clientRepo.FindByClientID(ctx, tokenCmd.ClientID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate client: %w", err)
	}
	if client == nil || client.ClientSecret != tokenCmd.ClientSecret {
		return nil, fmt.Errorf("invalid client credentials")
	}

	switch tokenCmd.GrantType {
	case "authorization_code":
		return h.handleAuthorizationCodeGrant(ctx, tokenCmd, client)
	case "client_credentials":
		return h.handleClientCredentialsGrant(ctx, tokenCmd, client)
	default:
		return nil, fmt.Errorf("unsupported grant_type: %s", tokenCmd.GrantType)
	}
}

func (h *TokenHandler) handleAuthorizationCodeGrant(ctx context.Context, cmd *TokenCommand, client *domain.Client) (interface{}, error) {
	// OAuth 2.1: Exact Redirect URI matching
	if client.RedirectURI != cmd.RedirectURI {
		return nil, fmt.Errorf("invalid redirect_uri")
	}

	// Find authorization code
	authCode, err := h.authCodeRepo.FindByCode(ctx, cmd.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve authorization code: %w", err)
	}
	if authCode == nil {
		return nil, fmt.Errorf("invalid authorization code")
	}

	// OAuth 2.1: Enforce PKCE Validation
	pkce := &domain.PKCE{
		CodeChallenge:       authCode.CodeChallenge,
		CodeChallengeMethod: authCode.CodeChallengeMethod,
	}
	if !pkce.ValidatePKCE(cmd.CodeVerifier) {
		return nil, fmt.Errorf("invalid code_verifier")
	}

	// Create access token
	scope := authCode.Scope
	if scope == "" {
		scope = "openid profile email"
	}
	
	accessToken := &domain.AccessToken{
		Token:     generateRandomToken(),
		UserID:    authCode.UserID,
		ClientID:  authCode.ClientID,
		Scope:     scope,
		ExpiresAt: getExpirationTime(),
	}

	if err := h.tokenRepo.Create(ctx, accessToken); err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	// OIDC: Generate ID Token if openid scope is present
	idToken := ""
	if strings.Contains(scope, "openid") {
		claims := map[string]interface{}{
			"sub": authCode.UserID,
			"aud": authCode.ClientID,
			"iss": "http://localhost:8080", // Should be config-driven
			"exp": time.Now().Add(1 * time.Hour).Unix(),
			"iat": time.Now().Unix(),
		}
		idToken, _ = h.jwtSigner.Sign(claims)
	}

	// Create refresh token
	refreshTokenStr := generateRandomToken()
	refreshToken := &domain.RefreshToken{
		Token:     refreshTokenStr,
		UserID:    authCode.UserID,
		ClientID:  authCode.ClientID,
		ExpiresAt: getExpirationTime().Add(30 * 24 * time.Hour), // 30 days
		Revoked:   false,
	}

	if err := h.refreshRepo.Create(ctx, refreshToken); err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	// Delete authorization code
	if err := h.authCodeRepo.Delete(ctx, cmd.Code); err != nil {
		return nil, fmt.Errorf("failed to delete authorization code: %w", err)
	}

	// Make Audit log entry
	auditEntry := &domain.Audit{
		ID:        generateRandomToken(),
		UserID:    authCode.UserID,
		ClientID:  authCode.ClientID,
		Action:    "ACCESS_TOKEN_ISSUED",
		Details:   "Access token and ID token issued via Authorization Code grant.",
		CreatedAt: time.Now(),
	}
	_ = h.auditRepo.Create(ctx, auditEntry)

	resp := map[string]interface{}{
		"access_token":  accessToken.Token,
		"refresh_token": refreshToken.Token,
		"token_type":    "Bearer",
		"expires_in":    3600,
	}
	if idToken != "" {
		resp["id_token"] = idToken
	}
	return resp, nil
}

func (h *TokenHandler) handleClientCredentialsGrant(ctx context.Context, cmd *TokenCommand, client *domain.Client) (interface{}, error) {
	// Create access token (system/client owned)
	accessToken := &domain.AccessToken{
		Token:     generateRandomToken(),
		UserID:    "system",
		ClientID:  client.ClientID,
		Scope:     cmd.Scope,
		ExpiresAt: getExpirationTime(),
	}

	if err := h.tokenRepo.Create(ctx, accessToken); err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	// Audit log
	auditEntry := &domain.Audit{
		ID:        generateRandomToken(),
		ClientID:  client.ClientID,
		Action:    "CLIENT_CREDENTIALS_TOKEN_ISSUED",
		Details:   "Access token issued via Client Credentials grant.",
		CreatedAt: time.Now(),
	}
	_ = h.auditRepo.Create(ctx, auditEntry)

	return map[string]interface{}{
		"access_token": accessToken.Token,
		"token_type":   "Bearer",
		"expires_in":   3600,
		"scope":        accessToken.Scope,
	}, nil
}
