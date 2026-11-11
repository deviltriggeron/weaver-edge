package dto

import "time"

type ContextResponse struct {
	Status   ContextStatus     `json:"status"`
	Messages []ContextMessage  `json:"messages"`
	Code     ContextDenialCode `json:"denial_code"`
}

type ContextMessage struct {
	ID        string    `json:"id"`
	ThreadID  string    `json:"thread_id"`
	Text      string    `json:"text"`
	SenderID  string    `json:"sender_id"`
	ChannelID string    `json:"channel_id"`
	CreatedAt time.Time `json:"created_at"`
}

type StatusResponse struct {
	Status string `json:"status"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type ContextStatus string

const (
	ContextStatusAllowed ContextStatus = "allowed"
	ContextStatusDenied  ContextStatus = "denied"
)

type ContextDenialCode string

const (
	ContextDenialNotAllowed  ContextDenialCode = "not_allowed"
	ContextDenialNotFound    ContextDenialCode = "not_found"
	ContextDenialUnsupported ContextDenialCode = "unsupported"
)
