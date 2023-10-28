package user_management

import (
	"context"

	"github.com/fahmifan/autograd/pkg/core"
	"github.com/fahmifan/autograd/pkg/core/auth"
	"github.com/fahmifan/autograd/pkg/dbmodel"
	"gorm.io/gorm"
)

type UserWriter struct{}

func (UserWriter) SaveUserWithPassword(ctx context.Context, tx *gorm.DB, user User, password auth.CipherPassword) error {
	model := dbmodel.User{
		Base: dbmodel.Base{
			ID:       user.ID,
			Metadata: core.NewModelMetadata(user.EntityMeta),
		},
		Name:     user.Name,
		Email:    user.Email,
		Password: string(password),
		Role:     string(user.Role),
	}

	return tx.Save(&model).Error
}
