package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base ..
type Base struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `sql:"index" json:"deleted_at"`
}

// BeforeCreate ..
func (base *Base) BeforeCreate(tx *gorm.DB) error {
	base.ID = uuid.New().String()
	return nil
}
