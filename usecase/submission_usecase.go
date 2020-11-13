package usecase

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/miun173/autograd/config"

	"github.com/miun173/autograd/repository"
)

// SubmissionUsecase ..
type SubmissionUsecase interface {
	Upload(fh *multipart.FileHeader) (string, error)
}

type submissionUsecase struct {
	submissionRepo repository.SubmissionRepository
}

// NewSubmissionUsecase ..
func NewSubmissionUsecase(submissionRepo repository.SubmissionRepository) SubmissionUsecase {
	return &submissionUsecase{
		submissionRepo: submissionRepo,
	}
}

func (e *submissionUsecase) Upload(fh *multipart.FileHeader) (string, error) {
	src, err := fh.Open()

	if err != nil {
		return "", err
	}

	defer src.Close()

	cwd, err := os.Getwd()

	if err != nil {
		return "", err
	}

	fileExt := filepath.Ext(fh.Filename)
	fileName := generateFileName(fh.Filename) + fileExt
	filePath := path.Join(cwd, "submission", fileName)
	dst, err := os.Create(filePath)

	if err != nil {
		return "", err
	}

	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	fileURL := config.BaseURL() + "/storage/" + fileName

	return fileURL, nil
}

func generateFileName(text string) string {
	h := md5.New()
	timestamp := fmt.Sprint(time.Now().Unix())

	h.Write([]byte(text + timestamp))

	return hex.EncodeToString(h.Sum(nil))
}
