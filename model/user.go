package model

import (
	"context"

	"github.com/fahmifan/autograd/pkg/core/auth"
)

// User ..
type User struct {
	Base
	Name     string
	Email    string
	Password string
	Role     auth.Role
}

// UserUsecase ..
type UserUsecase interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id string) (user *User, err error)
	FindByEmailAndPassword(ctx context.Context, email, password string) (*User, error)
}
