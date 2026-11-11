package repository

import (
	"context"
	"database/sql"
	"fmt"

	"weaver-edge/internal/domain"
	"weaver-edge/internal/infrastructure/db/mapper"
	"weaver-edge/internal/infrastructure/db/model"
	"weaver-edge/internal/ports"
)

type policyRepo struct {
	db *sql.DB
}

var _ ports.PolicyRepository = (*policyRepo)(nil)

func NewPolicyRepo(db *sql.DB) ports.PolicyRepository {
	return &policyRepo{
		db: db,
	}
}

func (p *policyRepo) Load(ctx context.Context) (*domain.PolicySnapshot, error) {
	q := `
		SELECT id, channel_id, thread_id, created_at
		FROM policies
	`

	rows, err := p.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("load allowed resources: %w", err)
	}
	defer rows.Close()

	snapshot := &domain.PolicySnapshot{
		Allowed: make([]domain.Policy, 0),
	}

	for rows.Next() {
		var policy model.PolicyModel

		err := rows.Scan(&policy.ID, &policy.ChannelID, &policy.ThreadID, &policy.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan allowed resource: %w", err)
		}

		snapshot.Allowed = append(snapshot.Allowed, mapper.PolicyModelToDomain(policy))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate allowed resources: %w", err)
	}

	return snapshot, nil
}
