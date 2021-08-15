package repository

import (
	"context"
	"errors"

	"github.com/fahmifan/autograd/model"
	"gorm.io/gorm"
)

// SessionRespository ..
type SessionRespository struct {
	db *gorm.DB
}

// NewSessionRepository ..
func NewSessionRepository(db *gorm.DB) *SessionRespository {
	return &SessionRespository{
		db: db,
	}
}

// CreateByUserID ..
func (s *SessionRespository) Create(ctx context.Context, sess *model.Session) error {
	return s.db.WithContext(ctx).Create(sess).Error
}

// FindByID ..
func (s *SessionRespository) FindByID(ctx context.Context, id string) (*model.Session, error) {
	sess := model.Session{}
	err := s.db.WithContext(ctx).Take(&sess, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &sess, err
}

// DeleteByID ..
func (s *SessionRespository) DeleteByID(ctx context.Context, id string) error {
	return s.db.WithContext(ctx).Delete(model.Session{}, "id = ?", id).Error
}

// DeleteAllByUserID ..
func (s *SessionRespository) DeleteAllByUserID(ctx context.Context, userID string) error {
	return s.db.WithContext(ctx).Delete(model.Session{}, "user_id = ?", userID).Error
}
