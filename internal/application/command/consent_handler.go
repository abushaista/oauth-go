package command

import (
	"context"
	"fmt"
	"github.com/abushaista/oauth-go/internal/domain"
	"github.com/abushaista/oauth-go/internal/infrastructure/security"
)

// ConsentHandler handles recording user consent
type ConsentHandler struct {
	consentRepo domain.ConsentRepository
}

// NewConsentHandler creates a new consent handler
func NewConsentHandler(consentRepo domain.ConsentRepository) *ConsentHandler {
	return &ConsentHandler{
		consentRepo: consentRepo,
	}
}

// Handle processes a consent command
func (h *ConsentHandler) Handle(ctx context.Context, cmd Command) (interface{}, error) {
	consentCmd, ok := cmd.(*ConsentCommand)
	if !ok {
		return nil, fmt.Errorf("invalid command type")
	}

	// Check if consent already exists
	existing, err := h.consentRepo.FindByUserAndClient(ctx, consentCmd.UserID, consentCmd.ClientID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing consent: %w", err)
	}

	if existing != nil {
		// Update existing consent
		existing.Scopes = consentCmd.Scopes
		if err := h.consentRepo.Update(ctx, existing); err != nil {
			return nil, fmt.Errorf("failed to update consent: %w", err)
		}
	} else {
		// Create new consent
		id, _ := security.GenerateRandomToken(16)
		consent := &domain.Consent{
			ID:       id,
			UserID:   consentCmd.UserID,
			ClientID: consentCmd.ClientID,
			Scopes:   consentCmd.Scopes,
		}

		if err := h.consentRepo.Create(ctx, consent); err != nil {
			return nil, fmt.Errorf("failed to store consent: %w", err)
		}
	}

	return nil, nil
}
