-- name: FindAllAssignmentsByCourseID :many
SELECT * FROM assignments WHERE id IN (
    SELECT assignment_id FROM rel_assignment_to_course WHERE course_id = @course_id
)
ORDER BY updated_at DESC
LIMIT @page_limit
OFFSET @page_offset;

-- name: FindCourseDetailForAssignmentByCourseID :one
SELECT id, "name", "description" FROM courses WHERE id = @id;