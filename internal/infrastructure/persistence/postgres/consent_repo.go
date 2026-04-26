package postgres

import (
	"context"
	"database/sql"
	"strings"

	"github.com/abushaista/oauth-go/internal/domain"
)

// ConsentRepository implements domain.ConsentRepository
type ConsentRepository struct {
	db *sql.DB
}

// NewConsentRepository creates a new consent repository
func NewConsentRepository(db *sql.DB) *ConsentRepository {
	return &ConsentRepository{db: db}
}

// Create stores a new consent record
func (r *ConsentRepository) Create(ctx context.Context, consent *domain.Consent) error {
	scopesStr := strings.Join(consent.Scopes, ",")
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO consents (id, user_id, client_id, scopes)
		 VALUES ($1, $2, $3, $4)`,
		consent.ID, consent.UserID, consent.ClientID, scopesStr,
	)
	return err
}

// FindByID retrieves a consent by ID
func (r *ConsentRepository) FindByID(ctx context.Context, id string) (*domain.Consent, error) {
	consent := &domain.Consent{}
	var scopesStr string

	err := r.db.QueryRowContext(
		ctx,
		"SELECT id, user_id, client_id, scopes FROM consents WHERE id = $1",
		id,
	).Scan(&consent.ID, &consent.UserID, &consent.ClientID, &scopesStr)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if scopesStr != "" {
		consent.Scopes = strings.Split(scopesStr, ",")
	}

	return consent, nil
}

// FindByUserAndClient retrieves consent for a specific user and client
func (r *ConsentRepository) FindByUserAndClient(ctx context.Context, userID, clientID string) (*domain.Consent, error) {
	consent := &domain.Consent{}
	var scopesStr string

	err := r.db.QueryRowContext(
		ctx,
		"SELECT id, user_id, client_id, scopes FROM consents WHERE user_id = $1 AND client_id = $2",
		userID, clientID,
	).Scan(&consent.ID, &consent.UserID, &consent.ClientID, &scopesStr)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if scopesStr != "" {
		consent.Scopes = strings.Split(scopesStr, ",")
	}

	return consent, nil
}

// Update modifies an existing consent
func (r *ConsentRepository) Update(ctx context.Context, consent *domain.Consent) error {
	scopesStr := strings.Join(consent.Scopes, ",")
	_, err := r.db.ExecContext(
		ctx,
		"UPDATE consents SET scopes = $1 WHERE id = $2",
		scopesStr, consent.ID,
	)
	return err
}

// Delete removes a consent record
func (r *ConsentRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(
		ctx,
		"DELETE FROM consents WHERE id = $1",
		id,
	)
	return err
}
