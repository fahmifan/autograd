
-- +migrate Up
CREATE TABLE IF NOT EXISTS "rel_assignment_to_courses" (
    "course_id" TEXT NOT NULL,
    "assignment_id" TEXT NOT NULL UNIQUE,
    FOREIGN KEY (course_id) REFERENCES courses(id),
    FOREIGN KEY (assignment_id) REFERENCES assignments(id)
);

CREATE INDEX idx_course_id ON "rel_assignment_to_courses" ("course_id");
CREATE INDEX idx_assignment_id ON "rel_assignment_to_courses" ("assignment_id");

-- +migrate Down

DROP TABLE IF EXISTS "rel_assignment_to_courses";
