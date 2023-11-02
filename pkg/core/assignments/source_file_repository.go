package assignments

import (
	"context"

	"github.com/fahmifan/autograd/pkg/dbmodel"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubmissionFileReader struct{}

func (SubmissionFileReader) FindByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (SubmissionFile, error) {
	file := dbmodel.File{}
	err := tx.WithContext(ctx).Where("id = ?", id).First(&file).Error
	return SubmissionFile{
		ID:  file.ID,
		URL: file.URL,
	}, err
}
