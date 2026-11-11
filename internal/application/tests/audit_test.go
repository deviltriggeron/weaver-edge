package tests

// TestNewAuditService	Сервис корректно создаётся
// TestAuditServiceSave	Успешное сохранение аудита
// TestAuditServiceSaveError	Репозиторий возвращает ошибку

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"weaver-edge/internal/application"
	"weaver-edge/internal/domain"
	"weaver-edge/internal/ports"
)

type fakeAuditRepository struct {
	event  domain.AuditEvent
	err    error
	called bool
}

var _ ports.AuditRepository = (*fakeAuditRepository)(nil)

func (f *fakeAuditRepository) Record(
	ctx context.Context,
	event domain.AuditEvent,
) error {
	f.called = true
	f.event = event
	return f.err
}

func TestAuditServiceSave(t *testing.T) {
	repo := &fakeAuditRepository{}
	service := application.NewAuditService(repo)

	before := time.Now().UTC()

	err := service.Save(context.Background(), domain.AuditEvent{})

	after := time.Now().UTC()

	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	if !repo.called {
		t.Fatal("repository Record() was not called")
	}

	if repo.event.ID == uuid.Nil {
		t.Error("event ID is nil, want generated UUID")
	}

	if repo.event.CreatedAt.IsZero() {
		t.Fatal("event CreatedAt is zero, want generated timestamp")
	}

	if repo.event.CreatedAt.Before(before) || repo.event.CreatedAt.After(after) {
		t.Errorf(
			"event CreatedAt = %v, want timestamp between %v and %v",
			repo.event.CreatedAt,
			before,
			after,
		)
	}

	if repo.event.CreatedAt.Location() != time.UTC {
		t.Errorf(
			"event CreatedAt location = %v, want UTC",
			repo.event.CreatedAt.Location(),
		)
	}
}

func TestAuditServiceSavePreservesExistingValues(t *testing.T) {
	wantID := uuid.New()
	wantCreatedAt := time.Date(
		2026, time.January, 15,
		10, 30, 0, 0,
		time.UTC,
	)

	event := domain.AuditEvent{
		ID:        wantID,
		CreatedAt: wantCreatedAt,
	}

	repo := &fakeAuditRepository{}
	service := application.NewAuditService(repo)

	err := service.Save(context.Background(), event)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	if repo.event.ID != wantID {
		t.Errorf("event ID = %v, want %v", repo.event.ID, wantID)
	}

	if !repo.event.CreatedAt.Equal(wantCreatedAt) {
		t.Errorf(
			"event CreatedAt = %v, want %v",
			repo.event.CreatedAt,
			wantCreatedAt,
		)
	}
}

func TestAuditServiceSaveError(t *testing.T) {
	wantErr := errors.New("database unavailable")

	repo := &fakeAuditRepository{
		err: wantErr,
	}
	service := application.NewAuditService(repo)

	err := service.Save(context.Background(), domain.AuditEvent{})

	if !errors.Is(err, wantErr) {
		t.Fatalf("Save() error = %v, want %v", err, wantErr)
	}

	if !repo.called {
		t.Fatal("repository Record() was not called")
	}
}
