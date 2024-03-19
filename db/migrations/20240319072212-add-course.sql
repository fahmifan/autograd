-- +migrate Up
CREATE TABLE IF NOT EXISTS "courses" (
    "id" TEXT PRIMARY KEY NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    "description" TEXT NOT NULL,
    "is_active" BOOLEAN NOT NULL,
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL,
    "deleted_at" TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "rel_course_users" (
    "course_id" TEXT NOT NULL,
    "user_id" TEXT NOT NULL,
    "user_type" VARCHAR(255) NOT NULL,
    FOREIGN KEY (course_id) REFERENCES courses(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- +migrate Down
DROP TABLE IF EXISTS "courses";
DROP TABLE IF EXISTS "rel_course_users";
