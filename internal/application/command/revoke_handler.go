package command

import (
	"context"
	"fmt"
	"time"
	"github.com/abushaista/oauth-go/internal/domain"
)

// RevokeCommand handles token revocation
type RevokeCommand struct {
	Token        string
	TokenType    string // "access_token" or "refresh_token"
	ClientID     string
	ClientSecret string
}

// RevokeHandler processes revocation requests
type RevokeHandler struct {
	tokenRepo   domain.AccessTokenRepository
	refreshRepo domain.RefreshTokenRepository
	clientRepo  domain.ClientRepository
	auditRepo   domain.AuditRepository
}

// NewRevokeHandler creates a new revoke handler
func NewRevokeHandler(
	tokenRepo domain.AccessTokenRepository,
	refreshRepo domain.RefreshTokenRepository,
	clientRepo domain.ClientRepository,
	auditRepo domain.AuditRepository,
) *RevokeHandler {
	return &RevokeHandler{
		tokenRepo:   tokenRepo,
		refreshRepo: refreshRepo,
		clientRepo:  clientRepo,
		auditRepo:   auditRepo,
	}
}

// Handle processes a revoke command
func (h *RevokeHandler) Handle(ctx context.Context, cmd Command) (interface{}, error) {
	revokeCmd, ok := cmd.(*RevokeCommand)
	if !ok {
		return nil, fmt.Errorf("invalid command type")
	}

	// 1. Authenticate client
	client, err := h.clientRepo.FindByClientID(ctx, revokeCmd.ClientID)
	if err != nil {
		return nil, fmt.Errorf("client lookup failed")
	}
	if client == nil || client.ClientSecret != revokeCmd.ClientSecret {
		return nil, fmt.Errorf("unauthorized client")
	}

	// 2. Locate token and revoke
	// RFC 7009: "If the token passed to the revocation endpoint has already been revoked..."
	// "...the revocation endpoint responds with an HTTP 200 OK."

	// Try access token
	if revokeCmd.TokenType == "" || revokeCmd.TokenType == "access_token" {
		_ = h.tokenRepo.Delete(ctx, revokeCmd.Token)
	}

	// Try refresh token
	if revokeCmd.TokenType == "" || revokeCmd.TokenType == "refresh_token" {
		_ = h.refreshRepo.Revoke(ctx, revokeCmd.Token)
	}

	// Audit log
	auditEntry := &domain.Audit{
		ID:        generateRandomToken(),
		ClientID:  revokeCmd.ClientID,
		Action:    "TOKEN_REVOKED",
		Details:   fmt.Sprintf("Token revocation requested for client %s", revokeCmd.ClientID),
		CreatedAt: time.Now(),
	}
	_ = h.auditRepo.Create(ctx, auditEntry)

	return nil, nil // Success
}
