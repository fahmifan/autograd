-- name: SaveUser :one
INSERT INTO users (id, "name", email, "password", "role", active, created_at, updated_at) 
VALUES (@id, @name, @email, @password, @role, @active, @created_at, @updated_at)
ON CONFLICT (id) DO UPDATE SET 
    "name" = EXCLUDED."name",
    email = EXCLUDED.email,
    "password" = EXCLUDED."password",
    "role" = EXCLUDED."role",
    active = EXCLUDED.active,
    created_at = EXCLUDED.created_at,
    updated_at = EXCLUDED.updated_at
RETURNING id;

-- name: SaveActivationToken :one
INSERT INTO activation_tokens (id, token, expired_at, created_at, updated_at) 
VALUES (@id, @token, @expired_at, @created_at, @updated_at)
ON CONFLICT (id) DO UPDATE SET 
    token = EXCLUDED.token,
    expired_at = EXCLUDED.expired_at,
    created_at = EXCLUDED.created_at,
    updated_at = EXCLUDED.updated_at
RETURNING id;

-- name: SaveRelUserToActivationToken :one
INSERT INTO rel_user_to_activation_tokens (user_id, activation_token_id)
VALUES (@user_id, @activation_token_id)
ON CONFLICT (user_id, activation_token_id) DO NOTHING
RETURNING user_id, activation_token_id;