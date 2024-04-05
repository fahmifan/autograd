-- name: SaveAssignment :exec
INSERT INTO assignments
(
    id, assigned_by, "name", "description", 
    "case_input_file_id", "case_output_file_id", 
    "deadline_at", "template",
    "created_at", "updated_at", "deleted_at"
)
VALUES
(
    @id, @assigned_by, @name, @description, 
    @case_input_file_id, @case_output_file_id,
    @deadline_at, @template,
    @created_at, @updated_at, @deleted_at
)
ON CONFLICT (id) DO UPDATE
SET
    assigned_by = @assigned_by, 
    "name" = @name, 
    "description" = @description,     
    "case_input_file_id" = @case_input_file_id, 
    "case_output_file_id" = @case_output_file_id, 
    "deadline_at" = @deadline_at, 
    "template" = @template,
    "created_at" = @created_at, 
    "updated_at" = @updated_at, 
    "deleted_at" = @deleted_at
;

-- name: SaveAssignmentCourse :exec
INSERT INTO rel_assignment_to_courses
(
    course_id, assignment_id
)
VALUES (@course_id, @assignment_id)
ON CONFLICT DO NOTHING;
