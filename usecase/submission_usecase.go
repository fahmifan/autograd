package usecase

import (
	"io"
	"math/rand"
	"mime/multipart"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/miun173/autograd/repository"
)

type SubmissionUsecase interface {
	Upload(fh *multipart.FileHeader) error
}

type submissionUsecase struct {
	submissionRepo repository.SubmissionRepository
}

func NewSubmissionUsecase(submissionRepo repository.SubmissionRepository) SubmissionUsecase {
	return &submissionUsecase{
		submissionRepo: submissionRepo,
	}
}

func (e *submissionUsecase) Upload(fh *multipart.FileHeader) error {
	src, err := fh.Open()

	if err != nil {
		return err
	}

	defer src.Close()

	cwd, err := os.Getwd()

	if err != nil {
		return err
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 16)
	randomString := generateRandomString(2)

	fileName := timestamp + "-" + randomString + "-" + fh.Filename
	filePath := path.Join(cwd, "submission", fileName)
	dst, err := os.Create(filePath)

	if err != nil {
		return err
	}

	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}

func generateRandomString(n int) string {
	letter := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, n)

	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}

	return string(b)
}
