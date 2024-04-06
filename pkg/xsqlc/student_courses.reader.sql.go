// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: student_courses.reader.sql

package xsqlc

import (
	"context"
)

const countAllStudentEnrolledCourses = `-- name: CountAllStudentEnrolledCourses :one
SELECT COUNT (id) FROM courses WHERE id IN (
    SELECT course_id FROM rel_course_users WHERE 
    user_id = $1 
    AND user_type = $2
)
`

type CountAllStudentEnrolledCoursesParams struct {
	UserID   string
	UserType string
}

func (q *Queries) CountAllStudentEnrolledCourses(ctx context.Context, arg CountAllStudentEnrolledCoursesParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, countAllStudentEnrolledCourses, arg.UserID, arg.UserType)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const findAllStudentEnrolledCourses = `-- name: FindAllStudentEnrolledCourses :many
SELECT id, "name", "description" FROM courses WHERE id IN (
    SELECT course_id FROM rel_course_users WHERE 
    user_id = $1 
    AND user_type = $2
)
ORDER BY id DESC
LIMIT $4
OFFSET $3
`

type FindAllStudentEnrolledCoursesParams struct {
	UserID     string
	UserType   string
	PageOffset int32
	PageLimit  int32
}

type FindAllStudentEnrolledCoursesRow struct {
	ID          string
	Name        string
	Description string
}

func (q *Queries) FindAllStudentEnrolledCourses(ctx context.Context, arg FindAllStudentEnrolledCoursesParams) ([]FindAllStudentEnrolledCoursesRow, error) {
	rows, err := q.db.QueryContext(ctx, findAllStudentEnrolledCourses,
		arg.UserID,
		arg.UserType,
		arg.PageOffset,
		arg.PageLimit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindAllStudentEnrolledCoursesRow
	for rows.Next() {
		var i FindAllStudentEnrolledCoursesRow
		if err := rows.Scan(&i.ID, &i.Name, &i.Description); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
