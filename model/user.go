package model

import (
	"context"
)

// User ..
type User struct {
	Base
	Name         string
	Email        string
	Password     string
	Role         Role
	RefreshToken string
}

// UserUsecase ..
type UserUsecase interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id string) (user *User, err error)
	FindByEmailAndPassword(ctx context.Context, email, password string) (*User, error)
}
