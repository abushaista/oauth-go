package domain

import "context"

type Consent struct {
	ID       string
	UserID   string
	ClientID string
	Scopes   []string
}

// ConsentRepository defines the interface for consent data operations
type ConsentRepository interface {
	// Create stores a new consent record
	Create(ctx context.Context, consent *Consent) error

	// FindByID retrieves a consent by ID
	FindByID(ctx context.Context, id string) (*Consent, error)

	// FindByUserAndClient retrieves consent for a specific user and client
	FindByUserAndClient(ctx context.Context, userID, clientID string) (*Consent, error)

	// Update modifies an existing consent
	Update(ctx context.Context, consent *Consent) error

	// Delete removes a consent record
	Delete(ctx context.Context, id string) error
}
