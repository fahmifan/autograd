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

// IsOwnedBy
func (s *Submission) IsOwnedBy(u User) bool {
	return s.SubmittedBy == u.ID
}

// SubmissionUsecase ..
type SubmissionUsecase interface {
	Create(ctx context.Context, submission *Submission) error
	DeleteByID(ctx context.Context, id string) (*Submission, error)
	FindByID(ctx context.Context, id string) (*Submission, error)
	FindByIDAndSubmitter(ctx context.Context, id, submitterID string) (*Submission, error)
	Update(ctx context.Context, submission *Submission) error
	FindAllByAssignmentID(ctx context.Context, cursor Cursor, assignmentID string) (submissions []*Submission, count int64, err error)
	UpdateGradeByID(ctx context.Context, id string, grade int64) error
	FindByAssignmentIDAndSubmitterID(ctx context.Context, assignmentID, submitterID string) ([]*Submission, error)
}
