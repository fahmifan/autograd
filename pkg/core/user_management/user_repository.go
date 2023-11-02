package user_management

import (
	"context"
	"fmt"

	"github.com/fahmifan/autograd/pkg/core"
	"github.com/fahmifan/autograd/pkg/core/auth"
	"github.com/fahmifan/autograd/pkg/dbmodel"
	"gorm.io/gorm"
)

type ManagedUserWriter struct{}

func (ManagedUserWriter) SaveUserWithPassword(ctx context.Context, tx *gorm.DB, user ManagedUser, password auth.CipherPassword) error {
	active := 0
	if user.Active {
		active = 1
	}

	model := dbmodel.User{
		Base: dbmodel.Base{
			ID:       user.ID,
			Metadata: core.NewModelMetadata(user.TimestampMetadata),
		},
		Name:     user.Name,
		Email:    user.Email,
		Password: string(password),
		Role:     string(user.Role),
		Active:   active,
	}

	return tx.Save(&model).Error
}

type ManagedUserReader struct{}

func (ManagedUserReader) FindUserByID(ctx context.Context, tx *gorm.DB, id string) (ManagedUser, error) {
	var model dbmodel.User
	if err := tx.First(&model, id).Error; err != nil {
		return ManagedUser{}, err
	}

	return ManagedUser{
		ID:                model.ID,
		Name:              model.Name,
		Email:             model.Email,
		Role:              auth.Role(model.Role),
		TimestampMetadata: core.TimestampMetaFromModel(model.Metadata),
		Active:            model.Active == 1,
	}, nil
}

type FindAllManagedUsersRequest struct {
	core.PaginationRequest
}

type FindAllManagedUsersResponse struct {
	Users []ManagedUser
	core.Pagination
}

func (ManagedUserReader) FindAll(ctx context.Context, tx *gorm.DB, req FindAllManagedUsersRequest) (res FindAllManagedUsersResponse, err error) {
	pagination := core.Pagination{
		Page:  req.Page,
		Limit: req.Limit,
	}

	var models []dbmodel.User
	if err := tx.Limit(int(req.Limit)).Offset(int(pagination.Offset())).Find(&models).Error; err != nil {
		return res, fmt.Errorf("find all: %w", err)
	}

	var count int64
	if err := tx.Model(&dbmodel.User{}).Count(&count).Error; err != nil {
		return res, fmt.Errorf("count: %w", err)
	}

	res.Users = make([]ManagedUser, len(models))
	for i, model := range models {
		res.Users[i] = ManagedUser{
			ID:                model.ID,
			Name:              model.Name,
			Email:             model.Email,
			Role:              auth.Role(model.Role),
			TimestampMetadata: core.TimestampMetaFromModel(model.Metadata),
			Active:            model.Active == 1,
		}
	}

	pagination.Total = int32(count)
	res.Pagination = pagination
	return res, nil
}
