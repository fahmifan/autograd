-- name: FindPendingByKey :many
SELECT * FROM outbox_items WHERE idempotent_key = $1 AND "status" = 'pending' LIMIT 1;
