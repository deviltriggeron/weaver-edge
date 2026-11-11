package repository

import (
	"context"
	"database/sql"
	"fmt"

	"weaver-edge/internal/domain"
	"weaver-edge/internal/infrastructure/db/mapper"
	"weaver-edge/internal/infrastructure/db/model"
	"weaver-edge/internal/ports"
)

type contextRepo struct {
	db *sql.DB
}

var _ ports.ContextRepository = (*contextRepo)(nil)

func NewContextRepository(db *sql.DB) ports.ContextRepository {
	return &contextRepo{
		db: db,
	}
}

func (r *contextRepo) SaveMessage(ctx context.Context, msg domain.Message) error {
	model := mapper.MsgDomainToModel(msg)
	q := `
		INSERT INTO messages(
			id, external_id, platform, channel_id, thread_id, sender_id, text, created_at, received_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.ExecContext(ctx, q, model.ID, model.ExternalID, model.Platform, model.ChannelID, model.ThreadID, model.SenderID, model.Text, model.CreatedAt, model.ReceivedAt)
	if err != nil {
		return fmt.Errorf("error save message: %w", err)
	}

	return nil
}

func (r *contextRepo) FetchContext(ctx context.Context, req domain.ContextRequest) ([]domain.Message, error) {
	if req.ChannelID == "" {
		return nil, fmt.Errorf("channel ID is required")
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}

	q := `
		SELECT 
			id, external_id, platform, channel_id, thread_id, sender_id, text, created_at, received_at
		FROM messages
		WHERE channel_id = $1
		AND (
			NULLIF($2, '') IS NULL
			OR thread_id = $2
		)
		ORDER BY created_at DESC, received_at DESC, id DESC
		LIMIT $3
	`

	rows, err := r.db.QueryContext(ctx, q, req.ChannelID, req.ThreadID, limit)
	if err != nil {
		return nil, fmt.Errorf("fetch context: %w", err)
	}
	defer rows.Close()

	messages := make([]domain.Message, 0, limit)

	for rows.Next() {
		var msg model.MessageModel

		err = rows.Scan(&msg.ID, &msg.ExternalID, &msg.Platform, &msg.ChannelID, &msg.ThreadID, &msg.SenderID, &msg.Text, &msg.CreatedAt, &msg.ReceivedAt)
		if err != nil {
			return nil, fmt.Errorf("scan context message: %w", err)
		}

		messages = append(messages, mapper.MsgModelToDomain(msg))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate context messages: %w", err)
	}

	return messages, nil
}
