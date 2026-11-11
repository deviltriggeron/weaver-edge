package model

import (
	"time"

	"github.com/google/uuid"
)

type MessageModel struct {
	ID         uuid.UUID `db:"id"`
	ExternalID string    `db:"external_id"`
	Platform   string    `db:"platform"`
	ThreadID   string    `db:"thread_id"`
	ChannelID  string    `db:"channel_id"`
	SenderID   string    `db:"sender_id"`
	Text       string    `db:"text"`
	CreatedAt  time.Time `db:"created_at"`
	ReceivedAt time.Time `db:"received_at"`
}
