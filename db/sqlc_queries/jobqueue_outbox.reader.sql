-- name: FindOutboxItemByByKey :one
SELECT * FROM outbox_items WHERE idempotent_key = $1 AND "status" = @status LIMIT 1;

-- name: FindOutboxItemByID :one
SELECT * FROM outbox_items WHERE id = $1;

-- name: FindAllOutboxItemIDsByStatus :many
SELECT id FROM outbox_items WHERE "status" = @status LIMIT @size_limit;