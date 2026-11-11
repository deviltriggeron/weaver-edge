package telegram

import (
	"context"

	"weaver-edge/internal/domain"
	"weaver-edge/internal/ports"
)

type telegramConnector struct {
	client telegramClient
}

var (
	_ ports.MessageReader = (*telegramConnector)(nil)
	_ ports.EventSource   = (*telegramConnector)(nil)
)

func NewTelegramConnector(client telegramClient) *telegramConnector {
	return &telegramConnector{
		client: client,
	}
}

func (c *telegramConnector) FetchMessage(ctx context.Context, id string) (*domain.Message, error) {
	return nil, nil
}

func (c *telegramConnector) FetchContext(ctx context.Context, req domain.ContextRequest) ([]domain.Message, error) {
	return nil, nil
}

func (c *telegramConnector) Receive(ctx context.Context) ([]domain.Message, error) {
	// updates, err := c.client.GetUpdates(ctx)
	// if err != nil {
	// 	return nil, err
	// }

	// var messages []domain.Message

	// fmt.Printf("Telegram: received %d updates\n", len(updates))

	// for _, update := range updates {
	// 	if update.Message == nil {
	// 		continue
	// 	}

	// 	messages = append(messages, dto.MapMessage(update.Message))
	// }

	// return messages, nil
	return nil, nil
}
