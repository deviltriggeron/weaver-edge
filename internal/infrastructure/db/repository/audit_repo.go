package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"weaver-edge/internal/domain"
	"weaver-edge/internal/ports"
)

type auditRepo struct {
	db *sql.DB
}

var _ ports.AuditRepository = (*auditRepo)(nil)

func NewAuditRepository(db *sql.DB) ports.AuditRepository {
	return &auditRepo{
		db: db,
	}
}

func (r *auditRepo) Record(ctx context.Context, event domain.AuditEvent) error {

	var decision []byte
	var err error

	if event.Decision != nil {
		decision, err = json.Marshal(event.Decision)
		if err != nil {
			return fmt.Errorf("marshal audit decision: %w", err)
		}
	}

	var parentEventID any
	if event.ParentEventID != nil {
		parentEventID = *event.ParentEventID
	}
	var messageID any
	if event.MessageID != nil {
		messageID = *event.MessageID
	}

	var errorMessage any
	if event.Error != "" {
		errorMessage = event.Error
	}

	q := `
		INSERT INTO audit_events (
			id, trace_id, parent_event_id, action, result, message_id, decision, error, created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err = r.db.ExecContext(ctx, q, event.ID, event.TraceID, parentEventID, event.Action, event.Result, messageID, decision, errorMessage, event.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert audit event: %w", err)
	}

	return nil
}
