package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// Assignment ..
type Assignment struct {
	ID                int64
	AssignedBy        int64
	Name              string
	Description       string
	CaseInputFileURL  string
	CaseOutputFileURL string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt
}

// AssignmentUsecase ..
type AssignmentUsecase interface {
	FindByID(ctx context.Context, id int64) (*Assignment, error)
}
