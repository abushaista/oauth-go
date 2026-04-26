package command

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/abushaista/oauth-go/internal/domain"
)

// AuthorizeHandler handles the authorization code flow
type AuthorizeHandler struct {
	authCodeRepo domain.AuthorizationCodeRepository
	clientRepo   domain.ClientRepository
	consentRepo  domain.ConsentRepository
	auditRepo    domain.AuditRepository
}

// NewAuthorizeHandler creates a new authorize handler
func NewAuthorizeHandler(
	authCodeRepo domain.AuthorizationCodeRepository,
	clientRepo domain.ClientRepository,
	consentRepo domain.ConsentRepository,
	auditRepo domain.AuditRepository,
) *AuthorizeHandler {
	return &AuthorizeHandler{
		authCodeRepo: authCodeRepo,
		clientRepo:   clientRepo,
		consentRepo:  consentRepo,
		auditRepo:    auditRepo,
	}
}

// Handle processes an authorization request
func (h *AuthorizeHandler) Handle(ctx context.Context, cmd Command) (interface{}, error) {
	authCmd, ok := cmd.(*AuthorizeCommand)
	if !ok {
		return nil, fmt.Errorf("invalid command type")
	}

	// Validate client
	client, err := h.clientRepo.FindByClientID(ctx, authCmd.ClientID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate client: %w", err)
	}
	if client == nil {
		return nil, fmt.Errorf("invalid client_id")
	}

	// Validate redirect URI
	if client.RedirectURI != authCmd.RedirectURI {
		return nil, fmt.Errorf("invalid redirect_uri")
	}

	// OAuth 2.1: Reject implicit flows
	if authCmd.ResponseType != "code" {
		return nil, fmt.Errorf("unsupported response_type: %s", authCmd.ResponseType)
	}

	// OAuth 2.1: Enforce PKCE
	if authCmd.CodeChallenge == "" {
		return nil, fmt.Errorf("code_challenge is required for OAuth 2.1 compliance")
	}

	// Check user consent
	consent, err := h.consentRepo.FindByUserAndClient(ctx, authCmd.UserID, authCmd.ClientID)
	if err != nil {
		return nil, fmt.Errorf("failed to check consent: %w", err)
	}

	// Verify all requested scopes are consented
	requestedScopes := strings.Split(authCmd.Scope, " ")
	hasAllConsent := consent != nil
	if hasAllConsent {
		consentedMap := make(map[string]bool)
		for _, s := range consent.Scopes {
			consentedMap[s] = true
		}
		for _, s := range requestedScopes {
			if s != "" && !consentedMap[s] {
				hasAllConsent = false
				break
			}
		}
	}

	if !hasAllConsent {
		// Consent needed - return to consent screen
		return map[string]interface{}{
			"require_consent": true,
			"client_id":       authCmd.ClientID,
			"redirect_uri":    authCmd.RedirectURI,
		}, nil
	}

	// Generate authorization code
	authCode := &domain.AuthorizationCode{
		Code:                generateRandomCode(),
		UserID:              authCmd.UserID,
		ClientID:            authCmd.ClientID,
		CodeChallenge:       authCmd.CodeChallenge,
		CodeChallengeMethod: authCmd.CodeChallengeMethod,
		ExpiresAt:           getExpirationTime(),
	}

	if err := h.authCodeRepo.Create(ctx, authCode); err != nil {
		return nil, fmt.Errorf("failed to create authorization code: %w", err)
	}

	// Make Audit log entry
	auditEntry := &domain.Audit{
		ID:        generateRandomToken(),
		UserID:    authCmd.UserID,
		ClientID:  authCmd.ClientID,
		Action:    "AUTHORIZATION_CODE_ISSUED",
		Details:   "Authorization code successfully issued to client.",
		IPAddress: "",
		CreatedAt: time.Now(),
	}
	_ = h.auditRepo.Create(ctx, auditEntry)

	return map[string]interface{}{
		"code":  authCode.Code,
		"state": authCmd.State,
	}, nil
}
