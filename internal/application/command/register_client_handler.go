package command

import (
	"context"
	"fmt"
	"time"
	"github.com/abushaista/oauth-go/internal/domain"
	"github.com/abushaista/oauth-go/internal/infrastructure/security"
)

// RegisterClientCommand handles dynamic client registration
type RegisterClientCommand struct {
	ClientName  string   `json:"client_name"`
	RedirectURIs []string `json:"redirect_uris"`
}

// RegisterClientHandler processes registration requests
type RegisterClientHandler struct {
	clientRepo domain.ClientRepository
	auditRepo  domain.AuditRepository
}

// NewRegisterClientHandler creates a new registration handler
func NewRegisterClientHandler(clientRepo domain.ClientRepository, auditRepo domain.AuditRepository) *RegisterClientHandler {
	return &RegisterClientHandler{
		clientRepo: clientRepo,
		auditRepo:  auditRepo,
	}
}

// Handle processes a registration command
func (h *RegisterClientHandler) Handle(ctx context.Context, cmd Command) (interface{}, error) {
	regCmd, ok := cmd.(*RegisterClientCommand)
	if !ok {
		return nil, fmt.Errorf("invalid command type")
	}

	clientID, _ := security.GenerateRandomToken(16)
	clientSecret, _ := security.GenerateRandomToken(32)
	id, _ := security.GenerateRandomToken(16)

	redirectURI := ""
	if len(regCmd.RedirectURIs) > 0 {
		redirectURI = regCmd.RedirectURIs[0]
	}

	client := &domain.Client{
		ID:           id,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
	}

	if err := h.clientRepo.Save(ctx, client); err != nil {
		return nil, fmt.Errorf("failed to save client: %w", err)
	}

	// Audit log
	auditEntry := &domain.Audit{
		ID:        generateRandomToken(),
		ClientID:  clientID,
		Action:    "CLIENT_REGISTERED",
		Details:   fmt.Sprintf("New client registered: %s", regCmd.ClientName),
		CreatedAt: time.Now(),
	}
	_ = h.auditRepo.Create(ctx, auditEntry)

	return client, nil
}
