package tests

// TestNewContextService	Сервис корректно создаётся
// TestContextServiceGetContext	Успешное получение контекста
// TestContextServiceGetContextError	Репозиторий возвращает ошибку
// TestContextServiceGetContextEmpty	Контекст отсутствует

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"weaver-edge/internal/application"
	"weaver-edge/internal/domain"
	"weaver-edge/internal/ports"
)

type fakeContextRepository struct {
	savedMessage domain.Message
	fetched      []domain.Message
	request      domain.ContextRequest

	saveErr  error
	fetchErr error

	saveCalled  bool
	fetchCalled bool
}

var _ ports.ContextRepository = (*fakeContextRepository)(nil)

func (f *fakeContextRepository) SaveMessage(
	ctx context.Context,
	msg domain.Message,
) error {
	f.saveCalled = true
	f.savedMessage = msg
	return f.saveErr
}

func (f *fakeContextRepository) FetchContext(
	ctx context.Context,
	req domain.ContextRequest,
) ([]domain.Message, error) {
	f.fetchCalled = true
	f.request = req
	return f.fetched, f.fetchErr
}

func TestContextServiceSaveMessage(t *testing.T) {
	message := domain.Message{
		ExternalID: "123",
		ChannelID:  "-100",
		ThreadID:   "4",
		Text:       "hello",
	}

	repo := &fakeContextRepository{}
	service := application.NewContextService(repo)

	err := service.SaveMessage(context.Background(), message)
	if err != nil {
		t.Fatalf("SaveMessage() error = %v", err)
	}

	if !repo.saveCalled {
		t.Fatal("repository SaveMessage() was not called")
	}

	if !reflect.DeepEqual(repo.savedMessage, message) {
		t.Errorf(
			"saved message = %+v, want %+v",
			repo.savedMessage,
			message,
		)
	}
}

func TestContextServiceSaveMessageError(t *testing.T) {
	wantErr := errors.New("database unavailable")

	repo := &fakeContextRepository{
		saveErr: wantErr,
	}
	service := application.NewContextService(repo)

	err := service.SaveMessage(context.Background(), domain.Message{})

	if !errors.Is(err, wantErr) {
		t.Fatalf("SaveMessage() error = %v, want %v", err, wantErr)
	}
}

func TestContextServiceGetContext(t *testing.T) {
	request := domain.ContextRequest{
		// Заполни поля, которые есть в твоём ContextRequest.
	}

	wantMessages := []domain.Message{
		{
			ExternalID: "123",
			ChannelID:  "-100",
			ThreadID:   "4",
			Text:       "first message",
		},
		{
			ExternalID: "124",
			ChannelID:  "-100",
			ThreadID:   "4",
			Text:       "second message",
		},
	}

	repo := &fakeContextRepository{
		fetched: wantMessages,
	}
	service := application.NewContextService(repo)

	gotMessages, err := service.GetContext(context.Background(), request)
	if err != nil {
		t.Fatalf("GetContext() error = %v", err)
	}

	if !repo.fetchCalled {
		t.Fatal("repository FetchContext() was not called")
	}

	if !reflect.DeepEqual(repo.request, request) {
		t.Errorf("request = %+v, want %+v", repo.request, request)
	}

	if !reflect.DeepEqual(gotMessages, wantMessages) {
		t.Errorf("messages = %+v, want %+v", gotMessages, wantMessages)
	}
}

func TestContextServiceGetContextError(t *testing.T) {
	wantErr := errors.New("database unavailable")

	repo := &fakeContextRepository{
		fetchErr: wantErr,
	}
	service := application.NewContextService(repo)

	_, err := service.GetContext(
		context.Background(),
		domain.ContextRequest{},
	)

	if !errors.Is(err, wantErr) {
		t.Fatalf("GetContext() error = %v, want %v", err, wantErr)
	}
}

func TestContextServiceGetContextEmpty(t *testing.T) {
	repo := &fakeContextRepository{
		fetched: []domain.Message{},
	}
	service := application.NewContextService(repo)

	got, err := service.GetContext(
		context.Background(),
		domain.ContextRequest{},
	)
	if err != nil {
		t.Fatalf("GetContext() error = %v", err)
	}

	if len(got) != 0 {
		t.Fatalf("GetContext() returned %d messages, want 0", len(got))
	}
}
