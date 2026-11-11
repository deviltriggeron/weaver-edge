package model

import (
	"time"

	"github.com/google/uuid"
)

type PendingReqModel struct {
	ID                uuid.UUID  `db:"id"`
	OriginalMessageID uuid.UUID  `db:"original_message_id"`
	ResponseMessageID *uuid.UUID `db:"response_message_id"`
	TraceID           uuid.UUID  `db:"trace_id"`
	Reason            string     `db:"reason"`
	Question          string     `db:"question"`
	Status            string     `db:"status"`
	CreatedAt         time.Time  `db:"created_at"`
	UpdatedAt         time.Time  `db:"updated_at"`
}
