-- name: FindAllCoursesByUser :many
SELECT * FROM courses
WHERE id IN (SELECT course_id FROM rel_course_users WHERE user_id = @user_id)
AND deleted_at IS NULL
LIMIT @page_limit
OFFSET @page_offset;

-- name: CountAllCoursesByUser :one
SELECT COUNT(*) FROM courses
WHERE id IN (SELECT course_id FROM rel_course_users WHERE user_id = @user_id)
AND deleted_at IS NULL;

-- name: FindCourseUserByID :one
SELECT id, name FROM users
WHERE id = @id;

-- name: FindCourseByID :one
SELECT * FROM courses
WHERE id = @id
AND deleted_at IS NULL;