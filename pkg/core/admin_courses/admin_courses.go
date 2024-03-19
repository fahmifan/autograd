package admin_courses

import (
	"errors"
	"strings"
	"time"

	"connectrpc.com/connect"
	"github.com/fahmifan/autograd/pkg/core"
	"github.com/fahmifan/ulids"
	"github.com/google/uuid"
)

var (
	ErrNameEmpty        = connect.NewError(connect.CodeInvalidArgument, errors.New("name cannot be empty"))
	ErrDescriptionEmpty = connect.NewError(connect.CodeInvalidArgument, errors.New("description cannot be empty"))
)

const CourseUserTypeAdmin = "admin"

type Course struct {
	ID          ulids.ULID
	Name        string
	Description string
	IsActive    bool
	Admin       Admin
	core.TimestampMetadata
}

type Admin struct {
	ID   uuid.UUID
	Name string
}

type CreateCourseRequest struct {
	Now         time.Time
	Admin       Admin
	NewID       ulids.ULID
	Name        string
	Description string
}

type UpdateCourseRequest struct {
	Now         time.Time
	Name        string
	Description string
	IsActive    bool
}

func CreateCourse(req CreateCourseRequest) (Course, error) {
	if len(strings.TrimSpace(req.Name)) <= 1 {
		return Course{}, ErrNameEmpty
	}

	if len(strings.TrimSpace(req.Description)) <= 1 {
		return Course{}, ErrDescriptionEmpty
	}

	return Course{
		ID:                req.NewID,
		Name:              req.Name,
		Description:       req.Description,
		IsActive:          false,
		Admin:             req.Admin,
		TimestampMetadata: core.NewTimestampMeta(req.Now),
	}, nil
}

func (course Course) Update(req UpdateCourseRequest) (Course, error) {
	if len(strings.TrimSpace(req.Name)) <= 1 {
		return Course{}, ErrNameEmpty
	}

	if len(strings.TrimSpace(req.Description)) <= 1 {
		return Course{}, ErrDescriptionEmpty
	}

	course.IsActive = req.IsActive
	course.Name = req.Name
	course.Description = req.Description
	course.UpdatedAt = req.Now

	return course, nil
}
