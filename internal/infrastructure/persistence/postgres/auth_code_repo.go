package postgres

import (
	"context"
	"database/sql"

	"github.com/abushaista/oauth-go/internal/domain"
)

// AuthorizationCodeRepository implements domain.AuthorizationCodeRepository
type AuthorizationCodeRepository struct {
	db *sql.DB
}

// NewAuthorizationCodeRepository creates a new authorization code repository
func NewAuthorizationCodeRepository(db *sql.DB) *AuthorizationCodeRepository {
	return &AuthorizationCodeRepository{db: db}
}

// Create stores a new authorization code
func (r *AuthorizationCodeRepository) Create(ctx context.Context, authCode *domain.AuthorizationCode) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO authorization_codes (code, user_id, client_id, code_challenge, code_challenge_method, scope, expires_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		authCode.Code, authCode.UserID, authCode.ClientID, authCode.CodeChallenge, authCode.CodeChallengeMethod, authCode.Scope, authCode.ExpiresAt,
	)
	return err
}

// FindByCode retrieves an authorization code by its code value
func (r *AuthorizationCodeRepository) FindByCode(ctx context.Context, code string) (*domain.AuthorizationCode, error) {
	authCode := &domain.AuthorizationCode{}
	err := r.db.QueryRowContext(
		ctx,
		`SELECT code, user_id, client_id, code_challenge, code_challenge_method, scope, expires_at
		 FROM authorization_codes WHERE code = $1 AND expires_at > NOW()`,
		code,
	).Scan(&authCode.Code, &authCode.UserID, &authCode.ClientID, &authCode.CodeChallenge, &authCode.CodeChallengeMethod, &authCode.Scope, &authCode.ExpiresAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return authCode, nil
}

// Delete removes an authorization code
func (r *AuthorizationCodeRepository) Delete(ctx context.Context, code string) error {
	_, err := r.db.ExecContext(
		ctx,
		"DELETE FROM authorization_codes WHERE code = $1",
		code,
	)
	return err
}
