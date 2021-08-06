CREATE TABLE IF NOT EXISTS "submissions" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "assignment_id" UUID REFERENCES assignments(id) NOT NULL,
    "submitted_by" UUID REFERENCES users(id) NOT NULL,
    "is_graded" BOOLEAN DEFAULT FALSE NOT NULL,
    "grade" INT NOT NULL DEFAULT 0,
    "feedback" TEXT DEFAULT '' NOT NULL,
    "file_url" TEXT DEFAULT '' NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL,
    "deleted_at" TIMESTAMPTZ
);