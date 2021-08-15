CREATE TABLE IF NOT EXISTS "sessions"(
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "user_id" UUID REFERENCES users(id) NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "expired_at" TIMESTAMPTZ NOT NULL
);
