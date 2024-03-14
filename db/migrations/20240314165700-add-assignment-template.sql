
-- +migrate Up
ALTER TABLE "assignments" ADD COLUMN template TEXT NOT NULL DEFAULT '';
-- +migrate Down
ALTER TABLE "assignments" DROP COLUMN template;
