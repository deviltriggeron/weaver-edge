package dto

import "time"

type SendMessageRequest struct {
	Message MessageRequest   `json:"message"`
	Context []MessageRequest `json:"context"`
}

type MessageRequest struct {
	ID         string    `json:"id"`
	ExternalID string    `json:"external_id"`
	Platform   string    `json:"platform"`
	ChannelID  string    `json:"channel_id"`
	ThreadID   string    `json:"thread_id"`
	SenderID   string    `json:"sender_id"`
	Text       string    `json:"text"`
	CreatedAt  time.Time `json:"created_at"`
	ReceivedAt time.Time `json:"received_at"`
}
