package repository

import (
	"context"

	"github.com/fahmifan/autograd/model"
	"github.com/fahmifan/autograd/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// UserRepository ..
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByID(ctx context.Context, id string) (*model.User, error)
}

type userRepo struct {
	db *gorm.DB
}

// NewUserRepository ..
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{
		db: db,
	}
}

// Create ..
func (u *userRepo) Create(ctx context.Context, user *model.User) (err error) {
	err = u.db.WithContext(ctx).Create(user).Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":  utils.Dump(ctx),
			"user": utils.Dump(user),
		}).Error(err)
	}

	return err
}

// FindByEmail find user by username. Upon not found will return nil, nil
func (u *userRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}
	err := u.db.WithContext(ctx).Where("email = ?", email).Take(user).Error
	switch err {
	case nil: // ignore
	case gorm.ErrRecordNotFound:
		return nil, nil
	default:
		logrus.WithFields(logrus.Fields{
			"ctx":   utils.Dump(ctx),
			"email": email,
		}).Error(err)
		return nil, err
	}

	return user, nil
}

// FindByID return nil, nil upon not found
func (u *userRepo) FindByID(ctx context.Context, id string) (*model.User, error) {
	user := &model.User{}
	err := u.db.WithContext(ctx).Where("id = ?", id).Take(user).Error
	switch err {
	case nil: // ignore
	case gorm.ErrRecordNotFound:
		return nil, nil
	default:
		logrus.WithFields(logrus.Fields{
			"ctx": utils.Dump(ctx),
			"id":  id,
		}).Error(err)
		return nil, err
	}

	return user, nil
}
