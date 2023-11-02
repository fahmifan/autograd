package assignments

import (
	"context"

	"github.com/fahmifan/autograd/pkg/dbmodel"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubmitterReader struct{}

func (SubmitterReader) FindByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (Submitter, error) {
	user := dbmodel.User{}
	err := tx.WithContext(ctx).Where("id = ?", id).First(&user).Error
	return Submitter{
		ID:     user.ID,
		Name:   user.Name,
		Active: user.Active == 1,
	}, err
}
