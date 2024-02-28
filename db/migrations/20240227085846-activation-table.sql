
-- +migrate Up
CREATE TABLE activation_tokens (
    id TEXT PRIMARY KEY NOT NULL,
    token TEXT NOT NULL,
    expired_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

CREATE TABLE rel_user_to_activation_tokens (
    user_id TEXT NOT NULL,
    activation_token_id TEXT NOT NULL,
    deleted_at TIMESTAMP,
    PRIMARY KEY (user_id, activation_token_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (activation_token_id) REFERENCES activation_tokens(id)
);

-- +migrate Down
DROP TABLE rel_user_to_activation_tokens;
DROP TABLE activation_tokens;
