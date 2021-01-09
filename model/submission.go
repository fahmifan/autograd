package model

import (
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
