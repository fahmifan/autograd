CREATE TABLE IF NOT EXISTS "assignments" (
    id BIGINT PRIMARY KEY,
    "name" TEXT NOT NULL,
    "description" TEXT NOT NULL,
    case_output_file TEXT NOT NULL,
    case_input_file TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ
);