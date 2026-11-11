-- +goose Up
CREATE TABLE pending_requests(
    id UUID PRIMARY KEY,
    original_message_id UUID NOT NULL REFERENCES messages(id),
    response_message_id UUID
        REFERENCES messages(id),
    trace_id UUID NOT NULL,
    reason TEXT NOT NULL,
    question TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_pending_requests_trace
ON pending_requests (trace_id);

CREATE INDEX idx_pending_requests_original_message
ON pending_requests (original_message_id);
-- +goose Down

DROP TABLE pending_requests;

