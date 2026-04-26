package postgres

import (
	"context"
	"database/sql"

	"github.com/abushaista/oauth-go/internal/domain"
)

// TokenRepository implements domain.AccessTokenRepository
type TokenRepository struct {
	db *sql.DB
}

// NewTokenRepository creates a new token repository
func NewTokenRepository(db *sql.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

// Create stores a new access token
func (r *TokenRepository) Create(ctx context.Context, token *domain.AccessToken) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO access_tokens (token, user_id, client_id, scope, expires_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		token.Token, token.UserID, token.ClientID, token.Scope, token.ExpiresAt,
	)
	return err
}

// FindByToken retrieves an access token by its token value
func (r *TokenRepository) FindByToken(ctx context.Context, token string) (*domain.AccessToken, error) {
	accessToken := &domain.AccessToken{}
	err := r.db.QueryRowContext(
		ctx,
		`SELECT token, user_id, client_id, scope, expires_at
		 FROM access_tokens WHERE token = $1 AND expires_at > NOW()`,
		token,
	).Scan(&accessToken.Token, &accessToken.UserID, &accessToken.ClientID, &accessToken.Scope, &accessToken.ExpiresAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return accessToken, nil
}

// Delete removes an access token
func (r *TokenRepository) Delete(ctx context.Context, token string) error {
	_, err := r.db.ExecContext(
		ctx,
		"DELETE FROM access_tokens WHERE token = $1",
		token,
	)
	return err
}
