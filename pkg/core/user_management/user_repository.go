package user_management

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fahmifan/autograd/pkg/core"
	"github.com/fahmifan/autograd/pkg/core/auth"
	"github.com/fahmifan/autograd/pkg/dbmodel"
	"github.com/google/uuid"
	"github.com/samber/lo"
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

	err := tx.Save(&model).Error
	if err != nil {
		return fmt.Errorf("SaveUserWithPassword: save user: %w", err)
	}

	tokenModel := dbmodel.ActivationToken{
		Base: dbmodel.Base{
			ID:       user.ActivationToken.ID,
			Metadata: core.NewModelMetadata(user.ActivationToken.TimestampMetadata),
		},
		Token:     user.ActivationToken.Token,
		ExpiredAt: user.ActivationToken.ExpiresAt,
	}

	err = tx.Save(&tokenModel).Error
	if err != nil {
		return fmt.Errorf("SaveUserWithPassword: save activation token: %w", err)
	}

	relModel := dbmodel.RelUserToActivationToken{
		UserID:            user.ID,
		ActivationTokenID: user.ActivationToken.ID,
	}
	err = tx.Save(&relModel).Error
	if err != nil {
		return fmt.Errorf("SaveUserWithPassword: save relation: %w", err)
	}

	return nil
}

func (ManagedUserWriter) SaveUser(ctx context.Context, tx *gorm.DB, user ManagedUser) error {
	active := 0
	if user.Active {
		active = 1
	}

	model := dbmodel.User{
		Base: dbmodel.Base{
			ID:       user.ID,
			Metadata: core.NewModelMetadata(user.TimestampMetadata),
		},
		Name:   user.Name,
		Email:  user.Email,
		Role:   string(user.Role),
		Active: active,
	}

	err := tx.Omit("password").Save(&model).Error
	if err != nil {
		return fmt.Errorf("SaveUser: save user: %w", err)
	}

	tokenModel := dbmodel.ActivationToken{
		Base: dbmodel.Base{
			ID:       user.ActivationToken.ID,
			Metadata: core.NewModelMetadata(user.ActivationToken.TimestampMetadata),
		},
		Token:     user.ActivationToken.Token,
		ExpiredAt: user.ActivationToken.ExpiresAt,
	}

	err = tx.Save(&tokenModel).Error
	if err != nil {
		return fmt.Errorf("SaveUserWithPassword: save activation token: %w", err)
	}

	return nil
}

type ManagedUserReader struct{}

func (ManagedUserReader) FindUserByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (ManagedUser, error) {
	var model dbmodel.User
	if err := tx.Take(&model, "id = ?", id).Error; err != nil {
		return ManagedUser{}, err
	}

	activationTokenModel, err := findActivationTokenByUserID(tx, id.String())
	if err != nil {
		return ManagedUser{}, fmt.Errorf("find activation token: %w", err)
	}

	return managedUserFromModel(model, activationTokenModel), nil
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
		res.Users[i] = managedUserFromModel(model, dbmodel.ActivationToken{})
	}

	userIDs := lo.Map(models, func(m dbmodel.User, _ int) uuid.UUID {
		return m.ID
	})
	findAllActivationTokenByUserIDs(tx, userIDs)

	pagination.Total = int32(count)
	res.Pagination = pagination
	return res, nil
}

func findActivationTokenByUserID(tx *gorm.DB, userID string) (dbmodel.ActivationToken, error) {
	var rel dbmodel.RelUserToActivationToken
	if err := tx.Take(&rel, "user_id = ?", userID).Error; err != nil {
		return dbmodel.ActivationToken{}, fmt.Errorf("find activation token: %w", err)
	}

	var model dbmodel.ActivationToken
	if err := tx.Take(&model, "id = ?", rel.ActivationTokenID).Error; err != nil {
		return dbmodel.ActivationToken{}, fmt.Errorf("find activation token: %w", err)
	}

	return model, nil
}

func findAllActivationTokenByUserIDs(tx *gorm.DB, userIDs []uuid.UUID) (dbmodel.ActivationToken, error) {
	var rel []dbmodel.RelUserToActivationToken
	if err := tx.Take(&rel, "user_id in ?", userIDs).Error; err != nil {
		return dbmodel.ActivationToken{}, fmt.Errorf("findAllActivationTokenByUserIDs: find relations: %w", err)
	}

	tokenIds := lo.Map(rel, func(r dbmodel.RelUserToActivationToken, _ int) uuid.UUID {
		return r.ActivationTokenID
	})

	now := time.Now()

	var model dbmodel.ActivationToken
	err := tx.Take(&model, "id in ? and ? < expires_at", tokenIds, now).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return dbmodel.ActivationToken{}, fmt.Errorf("findAllActivationTokenByUserIDs: find activation tokens: %w", err)
	}

	return model, nil
}

func managedUserFromModel(model dbmodel.User, activationTokenModel dbmodel.ActivationToken) ManagedUser {
	return ManagedUser{
		ID:                model.ID,
		Name:              model.Name,
		Email:             model.Email,
		Role:              auth.Role(model.Role),
		TimestampMetadata: core.TimestampMetaFromModel(model.Metadata),
		Active:            model.Active == 1,
		ActivationToken:   activationTokenFromModel(activationTokenModel),
	}
}

func activationTokenFromModel(model dbmodel.ActivationToken) ActivationToken {
	return ActivationToken{
		ID:    model.ID,
		Token: model.Token,
	}
}
