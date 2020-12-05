CREATE TABLE IF NOT EXISTS "submissions" (
    "id" BIGINT PRIMARY KEY,
    "assignment_id" BIGINT REFERENCES assignments(id) NOT NULL,
    "submitted_by" BIGINT REFERENCES users(id) NOT NULL,
    "is_graded" BOOLEAN DEFAULT FALSE NOT NULL,
    "grade" NUMERIC(3, 2) DEFAULT 0 NOT NULL,
    "feedback" TEXT DEFAULT '' NOT NULL,
    "file_url" TEXT DEFAULT '' NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL,
    "deleted_at" TIMESTAMPTZ
);