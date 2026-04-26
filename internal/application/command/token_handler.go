package command

import (
	"context"
	"fmt"
	"time"

	"github.com/abushaista/oauth-go/internal/domain"
)

// TokenHandler exchanges authorization codes and refresh tokens for access tokens
type TokenHandler struct {
	authCodeRepo domain.AuthorizationCodeRepository
	tokenRepo    domain.AccessTokenRepository
	refreshRepo  domain.RefreshTokenRepository
	clientRepo   domain.ClientRepository
	auditRepo    domain.AuditRepository
}

// NewTokenHandler creates a new token handler
func NewTokenHandler(
	authCodeRepo domain.AuthorizationCodeRepository,
	tokenRepo domain.AccessTokenRepository,
	refreshRepo domain.RefreshTokenRepository,
	clientRepo domain.ClientRepository,
	auditRepo domain.AuditRepository,
) *TokenHandler {
	return &TokenHandler{
		authCodeRepo: authCodeRepo,
		tokenRepo:    tokenRepo,
		refreshRepo:  refreshRepo,
		clientRepo:   clientRepo,
		auditRepo:    auditRepo,
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

	// Handle authorization code grant
	if tokenCmd.GrantType == "authorization_code" {
		return h.handleAuthorizationCodeGrant(ctx, tokenCmd, client)
	}

	return nil, fmt.Errorf("unsupported grant_type: %s", tokenCmd.GrantType)
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
	accessToken := &domain.AccessToken{
		Token:     generateRandomToken(),
		UserID:    authCode.UserID,
		ClientID:  authCode.ClientID,
		Scope:     "openid profile email",
		ExpiresAt: getExpirationTime(),
	}

	if err := h.tokenRepo.Create(ctx, accessToken); err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
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
		Details:   "Access token generated via Authorization Code grant.",
		IPAddress: "",
		CreatedAt: time.Now(),
	}
	_ = h.auditRepo.Create(ctx, auditEntry)

	return map[string]interface{}{
		"access_token":  accessToken.Token,
		"refresh_token": refreshToken.Token,
		"token_type":    "Bearer",
		"expires_in":    3600,
	}, nil
}
