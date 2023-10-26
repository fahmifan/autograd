CREATE TABLE "submissions" (
    "id" TEXT PRIMARY KEY NOT NULL,
    "assignment_id" TEXT NOT NULL,
    "submitted_by" TEXT REFERENCES users(id) NOT NULL,
    "is_graded" BOOLEAN DEFAULT FALSE NOT NULL,
    "grade" INT NOT NULL DEFAULT 0,
    "feedback" TEXT DEFAULT '' NOT NULL,
    "file_url" TEXT DEFAULT '' NOT NULL,
    "created_at" DATETIMET NOT NULL,
    "updated_at" DATETIMET NOT NULL,
    "deleted_at" DATETIMET,
    FOREIGN KEY REFERENCES assignments(id)
);