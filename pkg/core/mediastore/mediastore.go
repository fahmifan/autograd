package mediastore

import (
	"errors"
	"strings"
	"time"

	"github.com/fahmifan/autograd/pkg/core"
	"github.com/fahmifan/autograd/pkg/dbmodel"
	"github.com/google/uuid"
)

type MediaFileType dbmodel.FileType
type Extension string

type MediaFile struct {
	ID       uuid.UUID
	FileName string
	FilePath string
	FileType MediaFileType
	Ext      Extension
	URL      string
	core.TimestampMetadata
}

func (media MediaFile) GetInternalPath(baseURL string) string {
	strs := strings.Split(media.URL, baseURL)
	if len(strs) != 2 {
		return ""
	}

	return strs[1]
}

func ValidExtension(ext Extension) bool {
	switch ext {
	case ".txt", ".cpp":
		return true
	default:
		return false
	}
}

type CreateMediaRequest struct {
	NewID     uuid.UUID
	Now       time.Time
	FileName  string
	FilePath  string
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
		ID:                req.NewID,
		FileName:          req.FileName,
		FilePath:          req.FilePath,
		Ext:               req.Ext,
		FileType:          req.FileType,
		URL:               req.PublicURL,
		TimestampMetadata: core.NewTimestampMeta(req.Now),
	}, nil
}
