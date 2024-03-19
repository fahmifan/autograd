package admin_courses

import (
	"context"
	"fmt"

	"github.com/fahmifan/autograd/pkg/core"
	"github.com/fahmifan/autograd/pkg/xsqlc"
	"github.com/fahmifan/ulids"
	"github.com/google/uuid"
)

type AdminReader struct{}

func (r AdminReader) FindAdminByID(ctx context.Context, tx xsqlc.DBTX, id uuid.UUID) (Admin, error) {
	query := xsqlc.New(tx)
	res, err := query.FindCourseUserByID(ctx, id.String())
	if err != nil {
		return Admin{}, fmt.Errorf("FindCourseAdmin: %w", err)
	}

	return Admin{
		ID:   id,
		Name: res.Name,
	}, nil
}

type CourseReader struct{}

func (courseRd CourseReader) FindCourseByID(ctx context.Context, tx xsqlc.DBTX, id ulids.ULID) (Course, error) {
	query := xsqlc.New(tx)
	courseModel, err := query.FindCourseByID(ctx, id.String())
	if err != nil {
		return Course{}, fmt.Errorf("FindCourseByID: %w", err)
	}

	adminCourse, err := query.FindCourseUserByID(ctx, id.String())
	if err != nil {
		return Course{}, fmt.Errorf("FindCourseAdmin: %w", err)
	}

	return Course{
		ID:          id,
		Name:        courseModel.Name,
		Description: courseModel.Description,
		IsActive:    courseModel.IsActive,
		Admin: Admin{
			ID:   uuid.MustParse(adminCourse.ID),
			Name: adminCourse.Name,
		},
		TimestampMetadata: core.TimestampMetadata{
			CreatedAt: courseModel.CreatedAt,
			UpdatedAt: courseModel.UpdatedAt,
		},
	}, nil
}

type CourseWriter struct{}

func (courseWr CourseWriter) Save(ctx context.Context, tx xsqlc.DBTX, course *Course) error {
	query := xsqlc.New(tx)

	_, err := query.SaveCourse(ctx, xsqlc.SaveCourseParams{
		ID:          course.ID.String(),
		Name:        course.Name,
		Description: course.Description,
		IsActive:    course.IsActive,
		CreatedAt:   course.CreatedAt,
		UpdatedAt:   course.UpdatedAt,
	})

	if err != nil {
		return fmt.Errorf("SaveCourse: %w", err)
	}

	err = query.SaveRelCourseUser(ctx, xsqlc.SaveRelCourseUserParams{
		CourseID: course.ID.String(),
		UserID:   course.Admin.ID.String(),
		UserType: CourseUserTypeAdmin,
	})
	if err != nil {
		return fmt.Errorf("SaveRelCourseUser: %w", err)
	}

	return nil
}
