package usecase

import (
	"mime/multipart"
	"path"
	"path/filepath"

	"github.com/fahmifan/autograd/model"
	"github.com/fahmifan/autograd/utils"
	"github.com/sirupsen/logrus"
)

// MediaUsecase ..
type MediaUsecase struct {
	objectStorer model.ObjectStorer
	dstFolder    string
}

// NewMediaUsecase ..
func NewMediaUsecase(dstFolder string, uploader model.ObjectStorer) *MediaUsecase {
	return &MediaUsecase{
		objectStorer: uploader,
		dstFolder:    dstFolder,
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

	err = m.objectStorer.Store(dst, src)
	if err != nil {
		logrus.Error(err)
		return
	}

	return fileName, nil
}
