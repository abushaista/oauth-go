package postgres

import (
	"context"
	"database/sql"

	"github.com/abushaista/oauth-go/internal/domain"
)

// RefreshTokenRepository implements domain.RefreshTokenRepository
type RefreshTokenRepository struct {
	db *sql.DB
}

// NewRefreshTokenRepository creates a new refresh token repository
func NewRefreshTokenRepository(db *sql.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

// Create stores a new refresh token
func (r *RefreshTokenRepository) Create(ctx context.Context, token *domain.RefreshToken) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO refresh_tokens (token, user_id, client_id, revoked, expires_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		token.Token, token.UserID, token.ClientID, token.Revoked, token.ExpiresAt,
	)
	return err
}

// FindByToken retrieves a refresh token by its token value
func (r *RefreshTokenRepository) FindByToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	refreshToken := &domain.RefreshToken{}
	err := r.db.QueryRowContext(
		ctx,
		`SELECT token, user_id, client_id, revoked, expires_at
		 FROM refresh_tokens WHERE token = $1 AND expires_at > NOW()`,
		token,
	).Scan(&refreshToken.Token, &refreshToken.UserID, &refreshToken.ClientID, &refreshToken.Revoked, &refreshToken.ExpiresAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return refreshToken, nil
}

// Revoke marks a refresh token as revoked
func (r *RefreshTokenRepository) Revoke(ctx context.Context, token string) error {
	_, err := r.db.ExecContext(
		ctx,
		"UPDATE refresh_tokens SET revoked = TRUE WHERE token = $1",
		token,
	)
	return err
}

// Delete removes a refresh token
func (r *RefreshTokenRepository) Delete(ctx context.Context, token string) error {
	_, err := r.db.ExecContext(
		ctx,
		"DELETE FROM refresh_tokens WHERE token = $1",
		token,
	)
	return err
}
