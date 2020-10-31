package usecase

import (
	"context"
	"errors"

	"github.com/miun173/autograd/model"
	"github.com/miun173/autograd/repository"
	"github.com/miun173/autograd/utils"
	"github.com/sirupsen/logrus"
)

// UserUsecase ..
type UserUsecase interface {
	Create(ctx context.Context, user *model.User) error
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

// Create :nodoc:
func (u *userUsecase) Create(ctx context.Context, user *model.User) error {
	if user == nil {
		return errors.New("invalid arguments")
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
	err = u.userRepo.Create(ctx, user)
	if err != nil {
		logger.Error(err)
	}

	return err
}
