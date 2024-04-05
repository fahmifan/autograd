-- name: FindAllAssignmentsByCourseID :many
SELECT asg.*, rel.course_id FROM assignments asg
JOIN rel_assignment_to_courses rel ON asg.id = rel.assignment_id
WHERE rel.course_id = @course_id AND asg.deleted_at is NULL
ORDER BY updated_at DESC
LIMIT @page_limit
OFFSET @page_offset;

-- name: CountAllAssignmentsByCourse :one
SELECT COUNT(*) FROM assignments asg
JOIN rel_assignment_to_courses rel ON asg.id = rel.assignment_id
WHERE rel.course_id = @course_id AND asg.deleted_at is NULL;

-- name: FindAssignmentByID :one
SELECT asg.*, rel.course_id FROM assignments asg
JOIN rel_assignment_to_courses rel ON asg.id = rel.assignment_id
WHERE asg.id = @id AND asg.deleted_at is NULL;

-- name: FindCourseDetailForAssignmentByCourseID :one
SELECT id, "name", "description" FROM courses WHERE id = @id;