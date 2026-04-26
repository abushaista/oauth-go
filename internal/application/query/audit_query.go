package query

import (
	"context"
	"fmt"

	"github.com/abushaista/oauth-go/internal/domain"
)

// AuditQueryHandler retrieves audit event logs
type AuditQueryHandler struct {
	auditRepo domain.AuditRepository
}

// NewAuditQueryHandler creates a new audit query handler
func NewAuditQueryHandler(auditRepo domain.AuditRepository) *AuditQueryHandler {
	return &AuditQueryHandler{
		auditRepo: auditRepo,
	}
}

// Handle processes an audit query
func (h *AuditQueryHandler) Handle(ctx context.Context, q Query) (interface{}, error) {
	switch query := q.(type) {
	case *AuditUserQuery:
		limit := query.Limit
		if limit <= 0 || limit > 100 {
			limit = 50 // Default limit
		}
		return h.auditRepo.FindByUserID(ctx, query.UserID, limit)
	case *AuditClientQuery:
		limit := query.Limit
		if limit <= 0 || limit > 100 {
			limit = 50
		}
		return h.auditRepo.FindByClientID(ctx, query.ClientID, limit)
	default:
		return nil, fmt.Errorf("invalid query type for AuditQueryHandler")
	}
}
