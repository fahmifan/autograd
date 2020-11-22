package usecase

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/miun173/autograd/config"
	"github.com/miun173/autograd/model"
	"github.com/miun173/autograd/utils"
	"github.com/sirupsen/logrus"

	"github.com/miun173/autograd/repository"
)

// SubmissionUsecase ..
type SubmissionUsecase interface {
	Create(ctx context.Context, submission *model.Submission) error
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

func (s *submissionUsecase) Create(ctx context.Context, submission *model.Submission) error {
	if submission == nil {
		return errors.New("invalid arguments")
	}

	logger := logrus.WithFields(logrus.Fields{
		"ctx":        utils.Dump(ctx),
		"submission": utils.Dump(submission),
	})

	submission.ID = utils.GenerateID()
	submission.Feedback = ""
	submission.Grade = 0

	err := s.submissionRepo.Create(ctx, submission)
	if err != nil {
		logger.Error(err)
		return err
	}

	return err
}

func (s *submissionUsecase) Upload(fh *multipart.FileHeader) (string, error) {
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
