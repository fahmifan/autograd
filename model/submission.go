package model

import (
	"time"
)

// Submission ..
type Submission struct {
	ID           int64
	AssignmentID int64
	SubmittedBy  int64
	FileURL      string
	Grade        float64
	Feedback     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}
