
-- +migrate Up
CREATE INDEX outbox_items_idempotent_key ON outbox_items ("idempotent_key");

-- +migrate Down
DROP INDEX outbox_items_key;
