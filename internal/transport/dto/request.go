package dto

import "time"

type EventRequest struct {
	Message Message   `json:"message"`
	Context []Message `json:"context,omitempty"`
}

type Message struct {
	ID        string    `json:"id"`
	ThreadID  string    `json:"thread_id"`
	Text      string    `json:"text"`
	SenderID  string    `json:"sender_id"`
	ChannelID string    `json:"channel_id"`
	CreatedAt time.Time `json:"created_at"`
}

type ContextScope string

const (
	ContextScopeCurrentThread  ContextScope = "current_thread"
	ContextScopeCurrentChannel ContextScope = "current_channel"
)

type ContextRequest struct {
	ThreadID  string       `json:"thread_id"`
	ChannelID string       `json:"channel_id"`
	Scope     ContextScope `json:"scope"`
	Limit     int          `json:"limit"`
}
