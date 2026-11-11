package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"weaver-edge/config"
	"weaver-edge/internal/infrastructure/connectors/telegram/dto"
)

type telegramClient struct {
	baseURL    string
	botToken   string
	httpClient *http.Client
}

func NewTelegramClient(cfg config.Config) telegramClient {
	return telegramClient{
		baseURL:  "https://api.telegram.org",
		botToken: cfg.Connectors.Telegram,
		httpClient: &http.Client{
			Timeout: 35 * time.Second,
		},
	}
}

func (c *telegramClient) GetUpdates(ctx context.Context) ([]dto.Update, error) {
	url := fmt.Sprintf(
		"%s/bot%s/getUpdates?timeout=30",
		c.baseURL,
		c.botToken,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("telegram API returned status %d", resp.StatusCode)
	}

	var result dto.UpdateResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if !result.OK {
		return nil, fmt.Errorf("telegram API returned ok=false")
	}

	return result.Result, nil
}
