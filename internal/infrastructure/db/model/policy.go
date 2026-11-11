package model

import (
	"time"

	"github.com/google/uuid"
)

type PolicyModel struct {
	ID        uuid.UUID `db:"id"`
	ChannelID string    `db:"channel_id"`
	ThreadID  string    `db:"thread_id"`
	CreatedAt time.Time `db:"created_at"`
}
