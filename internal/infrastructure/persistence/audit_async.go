package persistence

import (
	"context"
	"log"

	"github.com/abushaista/oauth-go/internal/domain"
)

// AsyncAuditRepository wraps a domain.AuditRepository to handle Create operations asynchronously
type AsyncAuditRepository struct {
	repo  domain.AuditRepository
	queue chan *domain.Audit
}

// NewAsyncAuditRepository creates a new AsyncAuditRepository with a specified buffer size
func NewAsyncAuditRepository(repo domain.AuditRepository, bufferSize int) *AsyncAuditRepository {
	return &AsyncAuditRepository{
		repo:  repo,
		queue: make(chan *domain.Audit, bufferSize),
	}
}

// Create sends the audit log to the background worker
func (a *AsyncAuditRepository) Create(ctx context.Context, audit *domain.Audit) error {
	select {
	case a.queue <- audit:
		return nil
	default:
		// Queue is full - we log and drop to avoid blocking the main request flow
		log.Printf("Warning: Audit queue full, dropping audit log: %s", audit.Action)
		return nil // Return nil as we don't want to fail the main transaction due to audit logging
	}
}

// FindByUserID delegates to the underlying repository synchronously
func (a *AsyncAuditRepository) FindByUserID(ctx context.Context, userID string, limit int) ([]*domain.Audit, error) {
	return a.repo.FindByUserID(ctx, userID, limit)
}

// FindByClientID delegates to the underlying repository synchronously
func (a *AsyncAuditRepository) FindByClientID(ctx context.Context, clientID string, limit int) ([]*domain.Audit, error) {
	return a.repo.FindByClientID(ctx, clientID, limit)
}

// StartWorker starts the background worker that processes the audit queue
func (a *AsyncAuditRepository) StartWorker(ctx context.Context) {
	go func() {
		for {
			select {
			case audit := <-a.queue:
				// Use a fresh context for the DB operation to avoid it being cancelled
				// if the original request context is already done.
				err := a.repo.Create(context.Background(), audit)
				if err != nil {
					log.Printf("Error writing async audit log: %v", err)
				}
			case <-ctx.Done():
				// Handle graceful shutdown if needed, but for now just exit
				return
			}
		}
	}()
}
