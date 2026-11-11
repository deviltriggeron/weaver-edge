package application

import (
	"context"
	"time"

	"github.com/google/uuid"

	"weaver-edge/internal/domain"
	"weaver-edge/internal/ports"
)

type auditService struct {
	repo ports.AuditRepository
}

var _ ports.AuditService = (*auditService)(nil)

func NewAuditService(repo ports.AuditRepository) ports.AuditService {
	return &auditService{
		repo: repo,
	}
}

func (a *auditService) Save(ctx context.Context, event domain.AuditEvent) error {
	if event.ID == uuid.Nil {
		event.ID = uuid.New()
	}

	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now().UTC()
	}

	return a.repo.Record(ctx, event)
}
