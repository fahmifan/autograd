-- name: FindAllStudentsByCourse :many
SELECT u.id, u.name FROM users u
JOIN rel_course_users rel ON u.id = rel.user_id
WHERE rel.course_id = @course_id
ORDER BY rel.user_id DESC
LIMIT @page_limit
OFFSET @page_offset;

-- name: CountAllStudentsByCourse :one
SELECT COUNT(rel.user_id) FROM rel_course_users rel
WHERE rel.course_id = @course_id;