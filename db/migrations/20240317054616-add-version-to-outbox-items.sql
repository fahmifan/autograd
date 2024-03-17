
-- +migrate Up
ALTER TABLE outbox_items ADD COLUMN "version" INT NOT NULL DEFAULT 1;
-- +migrate Down
ALTER TABLE outbox_items DROP COLUMN "version";
