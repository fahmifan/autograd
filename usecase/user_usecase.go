package usecase

import (
	"context"
	"errors"

	"github.com/fahmifan/autograd/model"
	"github.com/fahmifan/autograd/repository"
	"github.com/fahmifan/autograd/utils"
	"github.com/sirupsen/logrus"
)

// UserUsecase ..
type UserUsecase interface {
	Create(ctx context.Context, user *model.User) error
	FindByID(ctx context.Context, id int64) (user *model.User, err error)
	FindByEmailAndPassword(ctx context.Context, email, password string) (*model.User, error)
}

type userUsecase struct {
	userRepo repository.UserRepository
}

// NewUserUsecase ..
func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

func (u *userUsecase) Create(ctx context.Context, user *model.User) error {
	if user == nil {
		return ErrInvalidArguments
	}

	logger := logrus.WithFields(logrus.Fields{
		"ctx":  utils.Dump(ctx),
		"user": utils.Dump(user),
	})

	oldUser, err := u.userRepo.FindByEmail(ctx, user.Email)
	if err != nil {
		logger.Error(err)
		return err
	}

	if oldUser != nil {
		return errors.New("user already exists")
	}

	user.ID = utils.GenerateID()
	user.Password, err = utils.HashPassword(user.Password)
	if err != nil {
		logrus.Error(err)
		return err
	}

	err = u.userRepo.Create(ctx, user)
	if err != nil {
		logger.Error(err)
	}

	return err
}

func (u *userUsecase) FindByID(ctx context.Context, id int64) (user *model.User, err error) {
	user, err = u.userRepo.FindByID(ctx, id)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx": utils.Dump(ctx),
			"id":  id,
		}).Error(err)
		return nil, err
	}

	if user == nil {
		return nil, ErrNotFound
	}

	return
}

func (u *userUsecase) FindByEmailAndPassword(ctx context.Context, email, plainPassword string) (user *model.User, err error) {
	user, err = u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":   utils.Dump(ctx),
			"email": email,
		}).Error(err)
		return nil, err
	}

	if user == nil {
		return nil, ErrNotFound
	}

	if !utils.CheckHashedPassword(plainPassword, user.Password) {
		return nil, ErrNotFound
	}

	return
}
