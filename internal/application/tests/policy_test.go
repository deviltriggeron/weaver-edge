package tests

import (
	"context"
	"errors"
	"testing"

	"weaver-edge/internal/application"
	"weaver-edge/internal/domain"
)

type fakePolicyRepository struct {
	snapshot *domain.PolicySnapshot
	err      error
}

func (m *fakePolicyRepository) Load(ctx context.Context) (*domain.PolicySnapshot, error) {
	return m.snapshot, m.err
}

func (m *fakePolicyRepository) Record(ctx context.Context, event domain.AuditEvent) error {
	return nil
}

// TestNewPolicyService_LoadError	Ошибка первоначальной загрузки
// TestNewPolicyService_NilSnapshot	Repository вернул nil
// TestPolicyService_Refresh	Успешная замена snapshot
// TestPolicyService_RefreshError	Старый snapshot сохраняется при ошибке
// TestPolicyService_RefreshNilSnapshot	Nil snapshot не заменяет старый

func TestPolicyEngineAllow(t *testing.T) {
	tests := []struct {
		name     string
		allowed  []domain.Policy
		resource domain.Resource
		want     bool
	}{
		{
			name: "allowed thread",
			allowed: []domain.Policy{
				{
					Resource: domain.Resource{
						ChannelID: "-1003995066662",
						ThreadID:  "5",
					},
				},
			},
			resource: domain.Resource{
				ChannelID: "-1003995066662",
				ThreadID:  "5",
			},
			want: true,
		},
		{
			name: "denied thread",
			allowed: []domain.Policy{
				{
					Resource: domain.Resource{
						ChannelID: "-1003995066662",
						ThreadID:  "5",
					},
				},
			},
			resource: domain.Resource{
				ChannelID: "-1003995066662",
				ThreadID:  "4",
			},
			want: false,
		},
		{
			name: "different channel",
			allowed: []domain.Policy{
				{
					Resource: domain.Resource{
						ChannelID: "-1003995066662",
						ThreadID:  "5",
					},
				},
			},
			resource: domain.Resource{
				ChannelID: "-5384174945",
				ThreadID:  "5",
			},
			want: false,
		},
		{
			name:    "no allowed resources",
			allowed: nil,
			resource: domain.Resource{
				ChannelID: "-1003995066662",
				ThreadID:  "5",
			},
			want: false,
		},
		{
			name: "second allowed resource",
			allowed: []domain.Policy{
				{
					Resource: domain.Resource{
						ChannelID: "-1003995066662",
						ThreadID:  "5",
					},
				},
				{
					Resource: domain.Resource{
						ChannelID: "-1003995066662",
						ThreadID:  "7",
					},
				},
			},
			resource: domain.Resource{
				ChannelID: "-1003995066662",
				ThreadID:  "7",
			},
			want: true,
		},
		{
			name: "whole channel allowed",
			allowed: []domain.Policy{
				{
					Resource: domain.Resource{
						ChannelID: "-1003995066662",
						ThreadID:  "",
					},
				},
			},
			resource: domain.Resource{
				ChannelID: "-1003995066662",
				ThreadID:  "5",
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakePolicyRepository{
				snapshot: &domain.PolicySnapshot{
					Allowed: tt.allowed,
				},
			}

			policy, err := application.NewPolicyService(context.Background(), repo)
			if err != nil {
				t.Fatalf("NewPolicyService() error = %v", err)
			}

			got := policy.Allow(context.Background(), tt.resource)

			if got != tt.want {
				t.Errorf("Allow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewPolicyServiceLoadError(t *testing.T) {
	expectedErr := errors.New("database error")

	repo := &fakePolicyRepository{
		err: expectedErr,
	}

	service, err := application.NewPolicyService(
		context.Background(),
		repo,
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if service != nil {
		t.Fatal("expected nil service")
	}

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected wrapped error %v, got %v", expectedErr, err)
	}
}

func TestNewPolicyServiceNilSnapshot(t *testing.T) {
	repo := &fakePolicyRepository{
		snapshot: nil,
	}

	service, err := application.NewPolicyService(
		context.Background(),
		repo,
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if service != nil {
		t.Fatal("expected nil service")
	}
}

func TestPolicyServiceRefresh(t *testing.T) {
	initialSnapshot := &domain.PolicySnapshot{
		Allowed: []domain.Policy{
			{
				Resource: domain.Resource{
					ChannelID: "channel-1",
					ThreadID:  "thread-1",
				},
			},
		},
	}

	updatedSnapshot := &domain.PolicySnapshot{
		Allowed: []domain.Policy{
			{
				Resource: domain.Resource{
					ChannelID: "channel-2",
					ThreadID:  "thread-2",
				},
			},
		},
	}

	repo := &fakePolicyRepository{
		snapshot: initialSnapshot,
	}

	ctx := context.Background()

	service, err := application.NewPolicyService(ctx, repo)
	if err != nil {
		t.Fatalf("NewPolicyService() error = %v", err)
	}

	repo.snapshot = updatedSnapshot

	if err := service.Refresh(ctx); err != nil {
		t.Fatalf("Refresh() error = %v", err)
	}

	if service.Allow(context.Background(), domain.Resource{
		ChannelID: "channel-2",
		ThreadID:  "thread-2",
	}) != true {
		t.Fatal("expected updated resource to be allowed")
	}

	if service.Allow(context.Background(), domain.Resource{
		ChannelID: "channel-1",
		ThreadID:  "thread-1",
	}) != false {
		t.Fatal("expected old resource to be denied")
	}
}

func TestPolicyServiceRefreshError(t *testing.T) {
	initialSnapshot := &domain.PolicySnapshot{
		Allowed: []domain.Policy{
			{
				Resource: domain.Resource{
					ChannelID: "channel-1",
					ThreadID:  "thread-1",
				},
			},
		},
	}

	repo := &fakePolicyRepository{
		snapshot: initialSnapshot,
	}

	ctx := context.Background()

	service, err := application.NewPolicyService(ctx, repo)
	if err != nil {
		t.Fatalf("NewPolicyService() error = %v", err)
	}

	expectedErr := errors.New("database unavailable")
	repo.err = expectedErr

	err = service.Refresh(context.Background())
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}

	allowed := service.Allow(context.Background(), domain.Resource{
		ChannelID: "channel-1",
		ThreadID:  "thread-1",
	})

	if !allowed {
		t.Fatal("old snapshot should remain active after refresh error")
	}
}

func TestPolicyServiceRefreshNilSnapshot(t *testing.T) {
	initialSnapshot := &domain.PolicySnapshot{
		Allowed: []domain.Policy{
			{
				Resource: domain.Resource{
					ChannelID: "channel-1",
					ThreadID:  "thread-1",
				},
			},
		},
	}

	repo := &fakePolicyRepository{
		snapshot: initialSnapshot,
	}

	service, err := application.NewPolicyService(
		context.Background(),
		repo,
	)
	if err != nil {
		t.Fatalf("NewPolicyService() error = %v", err)
	}

	repo.snapshot = nil

	if err := service.Refresh(context.Background()); err == nil {
		t.Fatal("expected error for nil snapshot")
	}

	if !service.Allow(context.Background(), domain.Resource{
		ChannelID: "channel-1",
		ThreadID:  "thread-1",
	}) {
		t.Fatal("old snapshot should remain active")
	}
}
