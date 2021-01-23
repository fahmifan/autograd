package usecase

import (
	"io"
	"mime/multipart"
	"path"
	"path/filepath"

	"github.com/miun173/autograd/utils"
	"github.com/sirupsen/logrus"
)

// Uploader ..
type Uploader interface {
	Upload(dst string, r io.Reader) error
}

// MediaUsecase ..
type MediaUsecase struct {
	uploader  Uploader
	dstFolder string
}

// NewMediaUsecase ..
func NewMediaUsecase(dstFolder string, uploader Uploader) *MediaUsecase {
	return &MediaUsecase{
		uploader:  uploader,
		dstFolder: dstFolder,
	}
}

// Upload ..
func (m *MediaUsecase) Upload(fileInfo *multipart.FileHeader) (name string, err error) {
	src, err := fileInfo.Open()
	if err != nil {
		logrus.Error(err)
		return
	}
	defer src.Close()

	ext := filepath.Ext(fileInfo.Filename)
	fileName := utils.GenerateUniqueString() + ext
	dst := path.Join(m.dstFolder, fileName)

	err = m.uploader.Upload(dst, src)
	if err != nil {
		logrus.Error(err)
		return
	}

	return fileName, nil
}
