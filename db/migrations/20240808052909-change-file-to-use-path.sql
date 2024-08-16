
-- +migrate Up
ALTER TABLE "files" ADD COLUMN "path" TEXT NOT NULL DEFAULT '';
-- +migrate Down
ALTER TABLE "files" DROP COLUMN "path";
