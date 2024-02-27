
-- +migrate Up
CREATE TABLE outbox_items (
    id TEXT PRIMARY KEY NOT NULL,
    idempotent_key TEXT NOT NULL,
    status TEXT NOT NULL,
    job_type TEXT NOT NULL,
    payload TEXT NOT NULL
);

-- +migrate Down
DROP TABLE outbox_items;
