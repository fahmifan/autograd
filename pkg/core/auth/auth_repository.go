package auth

import (
	"context"

	"github.com/fahmifan/autograd/pkg/dbmodel"
	"gorm.io/gorm"
)

type AuthReader struct{}

func (AuthReader) FindUserByEmail(ctx context.Context, tx *gorm.DB, email string) (authUser AuthUser, password CipherPassword, err error) {
	userModel := dbmodel.User{}
	err = tx.WithContext(ctx).Where("email = ?", email).Take(&userModel).Error
	if err != nil {
		return AuthUser{}, "", err
	}

	authUser = AuthUser{
		UserID: userModel.ID,
		Email:  userModel.Email,
		Name:   userModel.Name,
		Role:   Role(userModel.Role),
	}
	return authUser, CipherPassword(userModel.Password), nil
}
