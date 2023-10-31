package mediastore

import (
	"context"

	"github.com/fahmifan/autograd/pkg/dbmodel"
	"gorm.io/gorm"
)

type MediaFileWriter struct{}

func (MediaFileWriter) Create(ctx context.Context, tx *gorm.DB, mediaFile *MediaFile) error {
	model := dbmodel.File{
		Base: dbmodel.Base{
			ID:       mediaFile.ID,
			Metadata: mediaFile.ModelMetadata(),
		},
		Name: mediaFile.FileName,
		Type: dbmodel.FileType(mediaFile.FileType),
		Ext:  dbmodel.FileExt(mediaFile.Ext),
		URL:  mediaFile.URL,
	}

	return tx.Create(model).Error
}
