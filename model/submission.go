package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// Submission ..
type Submission struct {
	ID           int64
	AssignmentID int64
	SubmittedBy  int64
	FileURL      string
	Grade        int64
	Feedback     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt
}

// SubmissionUsecase ..
type SubmissionUsecase interface {
	FindByID(ctx context.Context, id int64) (*Submission, error)
	FindAllByAssignmentID(ctx context.Context, cursor Cursor, assignmentID int64) (submissions []*Submission, count int64, err error)
	UpdateGradeByID(ctx context.Context, id, grade int64) error
}
