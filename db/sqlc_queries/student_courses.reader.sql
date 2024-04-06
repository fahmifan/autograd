-- name: FindAllStudentEnrolledCourses :many
SELECT id, "name", "description" FROM courses WHERE id IN (
    SELECT course_id FROM rel_course_users WHERE 
    user_id = @user_id 
    AND user_type = @user_type
)
ORDER BY id DESC
LIMIT @page_limit
OFFSET @page_offset;

-- name: CountAllStudentEnrolledCourses :one
SELECT COUNT (id) FROM courses WHERE id IN (
    SELECT course_id FROM rel_course_users WHERE 
    user_id = @user_id 
    AND user_type = @user_type
);