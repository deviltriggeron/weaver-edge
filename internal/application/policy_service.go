package application

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"weaver-edge/internal/domain"
	"weaver-edge/internal/ports"
)

type policyService struct {
	snapshot *domain.PolicySnapshot
	mu       sync.RWMutex
	repo     ports.PolicyRepository
}

var _ ports.PolicyEngine = (*policyService)(nil)
var _ ports.PolicyRefresh = (*policyService)(nil)

// TODO:
// Add policy refresh mechanism:
// - reload snapshot from PolicyRepository;
// - replace current snapshot atomically;
// - define how policy changes are propagated across application instances.
func NewPolicyService(ctx context.Context, repo ports.PolicyRepository) (*policyService, error) {
	snapshot, err := repo.Load(ctx)
	if err != nil {
		return nil, fmt.Errorf("initialize policy: %w", err)
	}

	if snapshot == nil {
		return nil, fmt.Errorf("policy repository returned nil snapshot")
	}

	return &policyService{
		snapshot: snapshot,
		repo:     repo,
	}, nil
}

func (p *policyService) Allow(ctx context.Context, resource domain.Resource) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	for _, allowed := range p.snapshot.Allowed {
		if allowed.Resource.ChannelID != resource.ChannelID {
			continue
		}

		if allowed.Resource.ThreadID == "" || allowed.Resource.ThreadID == resource.ThreadID {
			return true
		}
	}

	return false
}

func (p *policyService) Refresh(ctx context.Context) error {
	snapshot, err := p.repo.Load(ctx)
	if err != nil {
		return fmt.Errorf("refresh policy: %w", err)
	}

	if snapshot == nil {
		return errors.New("policy repository returned nil snapshot")
	}

	p.mu.Lock()
	p.snapshot = snapshot
	p.mu.Unlock()

	return nil
}
