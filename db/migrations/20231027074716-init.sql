
-- +migrate Up

PRAGMA foreign_keys = ON;
CREATE TABLE "users" (
    "id" TEXT PRIMARY KEY NOT NULL,
    "name" TEXT NOT NULL,
    "email" TEXT NOT NULL,
    "password" TEXT NOT NULL,
    "role" INT NOT NULL,
    "active" INT NOT NULL DEFAULT 0,
    "created_at" DATETIMET NOT NULL,
    "updated_at" DATETIMET NOT NULL,
    "deleted_at" DATETIMET
);

CREATE TABLE "files" (
    "id" TEXT PRIMARY KEY NOT NULL,
    "name" TEXT NOT NULL,
    "type" TEXT NOT NULL,
    "ext" TEXT NOT NULL,
    "url" TEXT NOT NULL,
    "created_at" DATETIME NOT NULL,
    "updated_at" DATETIME NOT NULL,
    "deleted_at" DATETIME
);

CREATE TABLE "assignments" (
    "id" TEXT PRIMARY KEY NOT NULL,
    "assigned_by" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "description" TEXT NOT NULL,
    "case_output_file_id" TEXT NOT NULL,
    "case_input_file_id" TEXT NOT NULL,
    "created_at" DATETIME NOT NULL,
    "updated_at" DATETIME NOT NULL,
    "deleted_at" DATETIME,
    FOREIGN KEY (assigned_by) REFERENCES users(id)
);

CREATE TABLE "submissions" (
    "id" TEXT PRIMARY KEY NOT NULL,
    "assignment_id" TEXT NOT NULL,
    "submitted_by" TEXT REFERENCES users(id) NOT NULL,
    "is_graded" INT NOT NULL DEFAULT 0,
    "grade" INT NOT NULL DEFAULT 0,
    "feedback" TEXT DEFAULT '' NOT NULL,
    "file_url" TEXT DEFAULT '' NOT NULL,
    "created_at" DATETIMET NOT NULL,
    "updated_at" DATETIMET NOT NULL,
    "deleted_at" DATETIMET,
    FOREIGN KEY (submitted_by) REFERENCES users(id)
    FOREIGN KEY (assignment_id) REFERENCES assignments(id)
);

-- +migrate Down
DROP TABLE "submissions";
DROP TABLE "assignments";
DROP TABLE "users";

