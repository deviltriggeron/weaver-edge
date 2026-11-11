package telegram

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTelegramConnector_Receive(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		_, _ = w.Write([]byte(`{
			"ok": true,
			"result": [
				{
					"update_id": 100,
					"message": {
						"message_id": 42,
						"message_thread_id": 5,
						"from": {
							"id": 12345,
							"is_bot": false,
							"first_name": "Test",
							"username": "testuser"
						},
						"chat": {
							"id": -1003995066662,
							"type": "supergroup",
							"title": "Test Group"
						},
						"date": 1720000000,
						"text": "hello"
					}
				}
			]
		}`))
	}))
	defer server.Close()

	client := telegramClient{
		baseURL:    server.URL,
		botToken:   "test-token",
		httpClient: server.Client(),
	}

	connector := telegramConnector{
		client: client,
	}

	messages, err := connector.Receive(context.Background())
	if err != nil {
		t.Fatalf("Receive() error = %v", err)
	}

	if len(messages) != 1 {
		t.Fatalf("Receive() returned %d messages, want 1", len(messages))
	}

	message := messages[0]

	if message.ExternalID != "42" {
		t.Fatalf("ID = %q, want %q", message.ID, "42")
	}

	if message.ThreadID != "5" {
		t.Fatalf("ThreadID = %q, want %q", message.ThreadID, "5")
	}

	if message.ChannelID != "-1003995066662" {
		t.Fatalf("ChannelID = %q, want %q", message.ChannelID, "-1003995066662")
	}

	if message.SenderID != "12345" {
		t.Fatalf("SenderID = %q, want %q", message.SenderID, "12345")
	}

	if message.Text != "hello" {
		t.Fatalf("Text = %q, want %q", message.Text, "hello")
	}
}
