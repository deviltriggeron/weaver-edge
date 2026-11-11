package tests

import (
	"context"
	"log"
	"testing"

	"weaver-edge/config"
	"weaver-edge/internal/application"
	"weaver-edge/internal/domain"
	"weaver-edge/internal/infrastructure/db"
	"weaver-edge/internal/infrastructure/db/repository"
)

type fakeEventSource struct {
	messages []domain.Message
	err      error
}

func (f *fakeEventSource) Receive(ctx context.Context) ([]domain.Message, error) {
	return f.messages, f.err
}

type fakePolicy struct {
	deniedThread string
}

func (f *fakePolicy) Allow(ctx context.Context, resource domain.Resource) bool {
	return resource.ThreadID != f.deniedThread
}

func (f *fakePolicy) Refresh(ctx context.Context) error {
	return nil
}

type fakeCloudClient struct {
	messages []domain.Message
	err      error
}

func (f *fakeCloudClient) SendMessage(ctx context.Context, message domain.Message, context []domain.Message) (*domain.Decision, error) {
	if f.err != nil {
		return nil, f.err
	}

	f.messages = append(f.messages, message)
	return nil, nil
}

func TestCoreEventHandle(t *testing.T) {
	source := &fakeEventSource{}

	policy := &fakePolicy{
		deniedThread: "5",
	}

	core := &fakeCloudClient{}

	cfg := config.MustLoad()
	db, err := db.InitDB(cfg)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	contextRepo := repository.NewContextRepository(db)
	auditRepo := repository.NewAuditRepository(db)

	contextService := application.NewContextService(contextRepo)
	auditService := application.NewAuditService(auditRepo)

	handler := application.NewApplication(source, policy, core, contextService, auditService)

	msg := domain.Message{
		ExternalID: "1",
		Platform:   domain.Telegram,
		ChannelID:  "-100",
		ThreadID:   "4",
		Text:       "allowed",
	}

	err = handler.Handle(context.Background(), msg)
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}

	if len(core.messages) != 2 {
		t.Fatalf("sent %d messages, want 2", len(core.messages))
	}

	if core.messages[0].ExternalID != "1" {
		t.Fatalf("first message ID = %s, want 1", core.messages[0].ID)
	}

	if core.messages[1].ExternalID != "3" {
		t.Fatalf("second message ID = %s, want 3", core.messages[1].ID)
	}
}
