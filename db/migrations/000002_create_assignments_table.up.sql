CREATE TABLE "assignments" (
    "id" TEXT PRIMARY KEY DEFAULT,
    "assigned_by" TEXT REFERENCES users(id) NOT NULL,
    "name" TEXT NOT NULL,
    "description" TEXT NOT NULL,
    "case_output_file_url" TEXT NOT NULL,
    "case_input_file_url" TEXT NOT NULL,
    "created_at" DATETIME NOT NULL,
    "updated_at" DATETIME NOT NULL,
    "deleted_at" DATETIME
);