package model

import (
	"context"
	"time"
)

// User ..
type User struct {
	ID        int64
	Name      string
	Email     string
	Password  string
	Role      Role
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// UserUsecase ..
type UserUsecase interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id int64) (user *User, err error)
	FindByEmailAndPassword(ctx context.Context, email, password string) (*User, error)
}
