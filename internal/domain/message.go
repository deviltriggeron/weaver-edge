package domain

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID         uuid.UUID
	ExternalID string
	Platform   Platform
	ThreadID   string
	ChannelID  string
	SenderID   string
	Text       string
	CreatedAt  time.Time
	ReceivedAt time.Time
}
