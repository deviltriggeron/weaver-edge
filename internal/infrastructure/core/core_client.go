package core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"weaver-edge/config"
	"weaver-edge/internal/domain"
	"weaver-edge/internal/infrastructure/core/dto"
	"weaver-edge/internal/ports"
)

type coreClient struct {
	baseURL string
	token   string
	client  *http.Client
}

var _ ports.CoreClient = (*coreClient)(nil)

func NewCoreClient(client *http.Client, cfg config.Config) ports.CoreClient {
	return &coreClient{
		baseURL: cfg.Addr.CoreUrl,
		token:   cfg.Addr.CoreToken,
		client:  client,
	}
}

func (f *coreClient) SendMessage(ctx context.Context, message domain.Message, contextMessages []domain.Message) (*domain.Decision, error) {
	event := dto.MapSendMessageRequest(message, contextMessages)

	var body bytes.Buffer

	if err := json.NewEncoder(&body).Encode(event); err != nil {
		return nil, fmt.Errorf("encode message: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, f.baseURL+"/v1"+"/events", &body)
	if err != nil {
		return nil, fmt.Errorf("create core request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+f.token)

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send message to core: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("core returned HTTP status %s", resp.Status)
	}

	var responseDTO dto.EventResponse

	if err := json.NewDecoder(resp.Body).Decode(&responseDTO); err != nil {
		return nil, fmt.Errorf("decode core response: %w", err)
	}

	decision := dto.MapEventResponseToDecision(responseDTO)

	return &decision, nil
}
