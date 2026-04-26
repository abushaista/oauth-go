package command

import (
	"context"
	"fmt"
	"time"

	"github.com/abushaista/oauth-go/internal/domain"
)

// RefreshHandler handles token refresh requests
type RefreshHandler struct {
	refreshRepo domain.RefreshTokenRepository
	tokenRepo   domain.AccessTokenRepository
	clientRepo  domain.ClientRepository
	auditRepo   domain.AuditRepository
}

// NewRefreshHandler creates a new refresh handler
func NewRefreshHandler(
	refreshRepo domain.RefreshTokenRepository,
	tokenRepo domain.AccessTokenRepository,
	clientRepo domain.ClientRepository,
	auditRepo domain.AuditRepository,
) *RefreshHandler {
	return &RefreshHandler{
		refreshRepo: refreshRepo,
		tokenRepo:   tokenRepo,
		clientRepo:  clientRepo,
		auditRepo:   auditRepo,
	}
}

// Handle processes a refresh token request
func (h *RefreshHandler) Handle(ctx context.Context, cmd Command) (interface{}, error) {
	refreshCmd, ok := cmd.(*RefreshCommand)
	if !ok {
		return nil, fmt.Errorf("invalid command type")
	}

	// Validate client
	client, err := h.clientRepo.FindByClientID(ctx, refreshCmd.ClientID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate client: %w", err)
	}
	if client == nil || client.ClientSecret != refreshCmd.ClientSecret {
		return nil, fmt.Errorf("invalid client credentials")
	}

	// Find and validate refresh token
	refreshToken, err := h.refreshRepo.FindByToken(ctx, refreshCmd.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve refresh token: %w", err)
	}
	if refreshToken == nil || refreshToken.Revoked {
		return nil, fmt.Errorf("invalid or revoked refresh token")
	}

	// Generate new access token
	accessToken := &domain.AccessToken{
		Token:     generateRandomToken(),
		UserID:    refreshToken.UserID,
		ClientID:  refreshToken.ClientID,
		Scope:     refreshCmd.Scope,
		ExpiresAt: getExpirationTime(),
	}

	if err := h.tokenRepo.Create(ctx, accessToken); err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	// OAuth 2.1: Refresh Token Rotation
	// Revoke the old refresh token
	if err := h.refreshRepo.Revoke(ctx, refreshToken.Token); err != nil {
		return nil, fmt.Errorf("failed to revoke old refresh token: %w", err)
	}

	// Generate a new refresh token
	newRefreshTokenStr := generateRandomToken()
	newRefreshToken := &domain.RefreshToken{
		Token:     newRefreshTokenStr,
		UserID:    refreshToken.UserID,
		ClientID:  refreshToken.ClientID,
		ExpiresAt: getExpirationTime().Add(30 * 24 * time.Hour), // 30 days
		Revoked:   false,
	}

	if err := h.refreshRepo.Create(ctx, newRefreshToken); err != nil {
		return nil, fmt.Errorf("failed to create new refresh token: %w", err)
	}

	// Make Audit log entry
	auditEntry := &domain.Audit{
		ID:        generateRandomToken(),
		UserID:    refreshToken.UserID,
		ClientID:  refreshToken.ClientID,
		Action:    "TOKEN_REFRESHED",
		Details:   "Tokens fully rotated via Refresh exchange.",
		IPAddress: "",
		CreatedAt: time.Now(),
	}
	_ = h.auditRepo.Create(ctx, auditEntry)

	return map[string]interface{}{
		"access_token":  accessToken.Token,
		"refresh_token": newRefreshToken.Token,
		"token_type":    "Bearer",
		"expires_in":    3600,
	}, nil
}
