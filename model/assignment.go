package model

import (
	"context"
)

// Assignment ..
type Assignment struct {
	Base
	OwnedBy           string
	AssignedBy        string
	Name              string
	Description       string
	CaseInputFileURL  string
	CaseOutputFileURL string
}

func (a *Assignment) IsOwnedBy(u User) bool {
	return a.OwnedBy == u.ID
}

// AssignmentUsecase ..
type AssignmentUsecase interface {
	Create(ctx context.Context, assignment *Assignment) error
	DeleteByID(ctx context.Context, id string) (*Assignment, error)
	FindAll(ctx context.Context, cursor Cursor) (assignments []*Assignment, count int64, err error)
	FindByID(ctx context.Context, id string) (*Assignment, error)
	FindSubmissionsByID(ctx context.Context, cursor Cursor,
		assignmentID string) (submissions []*Submission, count int64, err error)
	Update(ctx context.Context, assignment *Assignment) error
}
