package model

import (
	"time"

	"github.com/google/uuid"
)

type AuditModel struct {
	ID            uuid.UUID  `db:"id"`
	TraceID       uuid.UUID  `db:"trace_id"`
	ParentEventID *uuid.UUID `db:"parent_event_id"`
	Action        string     `db:"action"`
	Result        string     `db:"result"`
	MessageID     *uuid.UUID `db:"message_id"`
	Decision      []byte     `db:"decision"`
	Error         *string    `db:"error"`
	CreatedAt     time.Time  `db:"created_at"`
}
