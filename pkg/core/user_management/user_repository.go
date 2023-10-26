package user_management

import (
	"context"

	"github.com/fahmifan/autograd/model"
	"github.com/fahmifan/autograd/pkg/core/auth"
	"gorm.io/gorm"
)

// User ..
type UserModel struct {
	model.Base
	Name     string
	Email    string
	Password string
	Role     auth.Role
}

func (UserModel) TableName() string {
	return "users"
}

type UserWriter struct{}

func (UserWriter) SaveUserWithPassword(ctx context.Context, tx *gorm.DB, user User, password string) error {
	model := UserModel{
		Base: model.Base{
			ID: user.ID.String(),
		},
		Name:     user.Name,
		Email:    user.Email,
		Password: password,
		Role:     user.Role,
	}
	return tx.Save(&model).Error
}
