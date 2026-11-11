package ports

import (
	"context"

	"weaver-edge/internal/domain"
)

type CoreClient interface {
	SendMessage(ctx context.Context, message domain.Message, contextMessages []domain.Message) (*domain.Decision, error)
}
