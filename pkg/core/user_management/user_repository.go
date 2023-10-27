package user_management

import (
	"context"

	"github.com/fahmifan/autograd/pkg/dbmodel"
	"gorm.io/gorm"
)

type UserWriter struct{}

func (UserWriter) SaveUserWithPassword(ctx context.Context, tx *gorm.DB, user User, password string) error {
	model := dbmodel.User{
		Base: dbmodel.Base{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			DeletedAt: gorm.DeletedAt(user.DeletedAt.NullTime),
		},
		Name:     user.Name,
		Email:    user.Email,
		Password: password,
		Role:     user.Role,
	}

	return tx.Save(&model).Error
}
