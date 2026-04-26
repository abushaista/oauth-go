package postgres

import (
	"context"
	"database/sql"

	"github.com/abushaista/oauth-go/internal/domain"
)

// AuditRepository implements domain.AuditRepository
type AuditRepository struct {
	db *sql.DB
}

// NewAuditRepository creates a new audit repository
func NewAuditRepository(db *sql.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

// Create stores a new audit log
func (r *AuditRepository) Create(ctx context.Context, audit *domain.Audit) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO audit_logs (id, user_id, client_id, action, details, ip_address, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		audit.ID, audit.UserID, audit.ClientID, audit.Action, audit.Details, audit.IPAddress, audit.CreatedAt,
	)
	return err
}

// FindByUserID retrieves audit logs for a specific user
func (r *AuditRepository) FindByUserID(ctx context.Context, userID string, limit int) ([]*domain.Audit, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, user_id, client_id, action, details, ip_address, created_at
		 FROM audit_logs WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2`,
		userID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var audits []*domain.Audit
	for rows.Next() {
		audit := &domain.Audit{}
		err := rows.Scan(&audit.ID, &audit.UserID, &audit.ClientID, &audit.Action, &audit.Details, &audit.IPAddress, &audit.CreatedAt)
		if err != nil {
			return nil, err
		}
		audits = append(audits, audit)
	}

	return audits, rows.Err()
}

// FindByClientID retrieves audit logs for a specific client
func (r *AuditRepository) FindByClientID(ctx context.Context, clientID string, limit int) ([]*domain.Audit, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, user_id, client_id, action, details, ip_address, created_at
		 FROM audit_logs WHERE client_id = $1 ORDER BY created_at DESC LIMIT $2`,
		clientID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var audits []*domain.Audit
	for rows.Next() {
		audit := &domain.Audit{}
		err := rows.Scan(&audit.ID, &audit.UserID, &audit.ClientID, &audit.Action, &audit.Details, &audit.IPAddress, &audit.CreatedAt)
		if err != nil {
			return nil, err
		}
		audits = append(audits, audit)
	}

	return audits, rows.Err()
}
