package postgres

import (
	"context"
	"database/sql"

	"github.com/abushaista/oauth-go/internal/domain"
)

// ClientRepository implements domain.ClientRepository
type ClientRepository struct {
	db *sql.DB
}

// NewClientRepository creates a new client repository
func NewClientRepository(db *sql.DB) *ClientRepository {
	return &ClientRepository{db: db}
}

// FindByClientID retrieves a client by client ID
func (r *ClientRepository) FindByClientID(ctx context.Context, clientID string) (*domain.Client, error) {
	client := &domain.Client{}
	err := r.db.QueryRowContext(
		ctx,
		"SELECT id, client_id, client_secret, redirect_uri FROM clients WHERE client_id = $1",
		clientID,
	).Scan(&client.ID, &client.ClientID, &client.ClientSecret, &client.RedirectURI)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return client, nil
}

// FindByID retrieves a client by ID
func (r *ClientRepository) FindByID(ctx context.Context, id string) (*domain.Client, error) {
	client := &domain.Client{}
	err := r.db.QueryRowContext(
		ctx,
		"SELECT id, client_id, client_secret, redirect_uri FROM clients WHERE id = $1",
		id,
	).Scan(&client.ID, &client.ClientID, &client.ClientSecret, &client.RedirectURI)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Save creates or updates a client
func (r *ClientRepository) Save(ctx context.Context, client *domain.Client) error {
	_, err := r.db.ExecContext(
		ctx,
		"INSERT INTO clients (id, client_id, client_secret, redirect_uri) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO UPDATE SET client_id = $2, client_secret = $3, redirect_uri = $4",
		client.ID, client.ClientID, client.ClientSecret, client.RedirectURI,
	)
	return err
}
