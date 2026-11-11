-- +goose Up
CREATE TABLE audit_events(
    id UUID PRIMARY KEY,
    trace_id UUID NOT NULL,
    parent_event_id UUID REFERENCES audit_events(id),
    action TEXT NOT NULL,
    result TEXT NOT NULL,
    message_id UUID REFERENCES messages(id) ON DELETE SET NULL,
    decision JSONB,
    error TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_events_trace
ON audit_events (trace_id, created_at);

CREATE INDEX idx_audit_events_parent
ON audit_events (parent_event_id);

CREATE INDEX idx_audit_events_message
ON audit_events (message_id);

-- +goose Down
DROP TABLE audit_events;
