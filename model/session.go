package model

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Session struct {
	ID        string
	UserID    string
	CreatedAt time.Time
	ExpiredAt time.Time
}

// BeforeCreate ..
func (s *Session) BeforeCreate(tx *gorm.DB) error {
	uuid, err := uuid.NewV4()
	if err != nil {
		return err
	}
	s.ID = uuid.String()
	return nil
}

type SessionRespository interface {
	FindByID(ctx context.Context, id string) (*Session, error)
	Create(ctx context.Context, sess *Session) error
	DeleteByID(ctx context.Context, id string) error
	DeleteAllByUserID(ctx context.Context, userID string) error
}
