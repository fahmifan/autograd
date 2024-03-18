-- name: CreateOutboxItem :one
INSERT INTO outbox_items (id, idempotent_key, "status", job_type, payload)
VALUES (@id, @idempotent_key, @status, @job_type, @payload)
RETURNING id, "version";

-- name: UpdateOutboxItem :one
UPDATE outbox_items
SET 
    "status" = @status,
    idempotent_key = @idempotent_key,
    job_type = @job_type,
    payload = @payload,
    "version" = "version" + 1
WHERE id = @id
    -- do optimistic locking
    AND "version" = "version"
    AND "version" = @version
RETURNING id, "version";