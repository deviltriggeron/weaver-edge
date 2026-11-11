package ports

import (
	"context"

	"weaver-edge/internal/domain"
)

type (
	ContextService interface {
		SaveMessage(ctx context.Context, msg domain.Message) error
		GetContext(ctx context.Context, request domain.ContextRequest) ([]domain.Message, error)
	}

	AuditService interface {
		Save(ctx context.Context, event domain.AuditEvent) error
	}

	PolicyEngine interface {
		Allow(ctx context.Context, resource domain.Resource) bool
	}

	PolicyRefresh interface {
		Refresh(ctx context.Context) error
	}

	EventSource interface {
		// Receive(ctx context.Context) ([]domain.Message, error)
	}
)
