package mediastore

import (
	"errors"

	"github.com/fahmifan/autograd/pkg/core"
	"github.com/fahmifan/autograd/pkg/dbmodel"
	"github.com/google/uuid"
)

type MediaFileType dbmodel.FileType
type Extension string

type MediaFile struct {
	ID       uuid.UUID
	FileName string
	FileType MediaFileType
	Ext      Extension
	URL      string
	core.TimestampMetadata
}

func ValidExtension(ext Extension) bool {
	switch ext {
	case ".txt":
		return true
	default:
		return false
	}
}

type CreateMediaRequest struct {
	NewID     uuid.UUID
	FileName  string
	FileType  MediaFileType
	Ext       Extension
	PublicURL string
}

func CreateMediaFile(req CreateMediaRequest) (MediaFile, error) {
	if !ValidExtension(req.Ext) {
		return MediaFile{}, errors.New("invalid extension")
	}

	if req.FileName == "" {
		return MediaFile{}, errors.New("invalid file name")
	}

	return MediaFile{
		ID:       req.NewID,
		FileName: req.FileName,
		Ext:      req.Ext,
		FileType: req.FileType,
	}, nil
}
