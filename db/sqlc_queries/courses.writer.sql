-- name: SaveCourse :one
INSERT INTO courses (
    id, 
    "name", 
    "description", 
    is_active, 
    created_at, 
    updated_at
)
VALUES (
    @id,
    @name,
    @description,
    @is_active,
    @created_at,
    @updated_at
)
ON CONFLICT (id) DO UPDATE
SET
    "name" = @name,
    "description" = @description,
    is_active = @is_active,
    updated_at = @updated_at
RETURNING id;

-- name: SaveRelCourseUser :exec
INSERT INTO rel_course_users (
    course_id, 
    user_id, 
    user_type
)
VALUES (
    @course_id,
    @user_id,
    @user_type
);