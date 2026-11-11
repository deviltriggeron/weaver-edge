package application

import (
	"context"

	"weaver-edge/internal/domain"
	"weaver-edge/internal/ports"
)

type contextService struct {
	db ports.ContextRepository
}

var _ ports.ContextService = (*contextService)(nil)

func NewContextService(db ports.ContextRepository) ports.ContextService {
	return &contextService{
		db: db,
	}
}

func (c *contextService) SaveMessage(ctx context.Context, msg domain.Message) error {
	return c.db.SaveMessage(ctx, msg)
}

func (c *contextService) GetContext(ctx context.Context, request domain.ContextRequest) ([]domain.Message, error) {
	return c.db.FetchContext(ctx, request)
}
