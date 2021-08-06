package model

import (
	"context"
)

// Submission ..
type Submission struct {
	Base
	AssignmentID string
	SubmittedBy  string
	FileURL      string
	Grade        int64
	Feedback     string
}

// SubmissionUsecase ..
type SubmissionUsecase interface {
	Create(ctx context.Context, submission *Submission) error
	DeleteByID(ctx context.Context, id string) (*Submission, error)
	FindByID(ctx context.Context, id string) (*Submission, error)
	Update(ctx context.Context, submission *Submission) error
	FindAllByAssignmentID(ctx context.Context, cursor Cursor, assignmentID string) (submissions []*Submission, count int64, err error)
	UpdateGradeByID(ctx context.Context, id string, grade int64) error
}
