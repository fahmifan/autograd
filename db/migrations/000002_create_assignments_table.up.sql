CREATE TABLE IF NOT EXISTS "assignments" (
    "id" BIGINT PRIMARY KEY,
    "assigned_by" BIGINT REFERENCES users(id) NOT NULL,
    "name" TEXT NOT NULL,
    "description" TEXT NOT NULL,
    "case_output_file_url" TEXT NOT NULL,
    "case_input_file_url" TEXT NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL,
    "deleted_at" TIMESTAMPTZ
);