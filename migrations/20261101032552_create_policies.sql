-- +goose Up
CREATE TABLE policies(
    id UUID PRIMARY KEY,
    channel_id TEXT NOT NULL,
    thread_id TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (channel_id, thread_id)
);

-- +goose Down
DROP TABLE policies;