package usecase

import (
	"fmt"
	"io"
	"mime/multipart"
	"path"
	"path/filepath"

	"github.com/miun173/autograd/config"
	"github.com/miun173/autograd/utils"
	"github.com/sirupsen/logrus"
)

// Uploader ..
type Uploader interface {
	Upload(dst string, r io.Reader) error
}

// MediaUsecase ..
type MediaUsecase struct {
	uploader Uploader
}

// NewMediaUsecase ..
func NewMediaUsecase(uploader Uploader) *MediaUsecase {
	return &MediaUsecase{
		uploader: uploader,
	}
}

// Upload ..
func (m *MediaUsecase) Upload(fileInfo *multipart.FileHeader) (pubURL string, err error) {
	src, err := fileInfo.Open()
	if err != nil {
		logrus.Error(err)
		return
	}
	defer src.Close()

	ext := filepath.Ext(fileInfo.Filename)
	fileName := utils.GenerateUniqueString() + ext
	dst := path.Join("media", fileName)

	err = m.uploader.Upload(dst, src)
	if err != nil {
		logrus.Error(err)
		return
	}

	publicURL := fmt.Sprintf("%s/%s", config.BaseURL(), dst)
	return publicURL, nil
}
