package ports

import (
	"context"

	"weaver-edge/internal/domain"
)

type (
	ContextRepository interface {
		SaveMessage(ctx context.Context, msg domain.Message) error
		FetchContext(ctx context.Context, req domain.ContextRequest) ([]domain.Message, error)
	}

	PolicyRepository interface {
		Load(ctx context.Context) (*domain.PolicySnapshot, error)
	}

	AuditRepository interface {
		Record(ctx context.Context, event domain.AuditEvent) error
	}

	// think about whether this is necessary
	MessageReader interface {
		FetchMessage(ctx context.Context, id string) (*domain.Message, error)
		FetchContext(ctx context.Context, req domain.ContextRequest) ([]domain.Message, error)
	}
)
