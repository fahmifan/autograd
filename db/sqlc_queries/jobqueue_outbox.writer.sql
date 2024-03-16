-- name: CreateOutboxItem :one
INSERT INTO outbox_items (id, idempotent_key, "status", job_type, payload)
VALUES (@id, @idempotent_key, @status, @job_type, @payload)
RETURNING id;