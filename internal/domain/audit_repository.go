package domain

import "context"

// AuditRepository defines the interface for audit log operations
type AuditRepository interface {
	// Create stores a new audit log
	Create(ctx context.Context, audit *Audit) error

	// FindByUserID retrieves audit logs for a specific user
	FindByUserID(ctx context.Context, userID string, limit int) ([]*Audit, error)

	// FindByClientID retrieves audit logs for a specific client
	FindByClientID(ctx context.Context, clientID string, limit int) ([]*Audit, error)
}
