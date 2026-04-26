package domain

import "context"

type Client struct {
	ID           string
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

// ClientRepository defines the interface for client data operations
type ClientRepository interface {
	// FindByClientID retrieves a client by client ID
	FindByClientID(ctx context.Context, clientID string) (*Client, error)

	// FindByID retrieves a client by ID
	FindByID(ctx context.Context, id string) (*Client, error)

	// Save creates or updates a client
	Save(ctx context.Context, client *Client) error
}
