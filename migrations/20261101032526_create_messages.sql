-- +goose Up
CREATE TABLE messages(
    id UUID PRIMARY KEY,
    platform TEXT NOT NULL,
    external_id TEXT NOT NULL,
    channel_id TEXT NOT NULL,
    thread_id TEXT NOT NULL DEFAULT '',
    sender_id TEXT NOT NULL,
    text TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    received_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (platform, external_id)
);

CREATE INDEX idx_messages_channel_thread_created
ON messages (channel_id, thread_id, created_at);

-- +goose Down
DROP TABLE messages;