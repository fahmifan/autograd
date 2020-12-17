package usecase

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path"
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
	Upload(ctx context.Context, sourceCode string) (string, error)
	FindAllByAssignmentID(ctx context.Context, cursor model.Cursor, assignmentID int64) (submissions []*model.Submission, count int64, err error)
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
	err := s.submissionRepo.Create(ctx, submission)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (s *submissionUsecase) Upload(ctx context.Context, sourceCode string) (string, error) {
	if sourceCode == "" {
		return "", errors.New("invalid arguments")
	}

	logger := logrus.WithFields(logrus.Fields{
		"ctx":        utils.Dump(ctx),
		"sourceCode": utils.Dump(sourceCode),
	})

	cwd, err := os.Getwd()
	if err != nil {
		logger.Error(err)
		return "", err
	}

	fileName := generateFileName() + ".cpp"
	filePath := path.Join(cwd, "submission", fileName)
	file, err := os.Create(filePath)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	file.WriteString(sourceCode)
	defer file.Close()

	fileURL := config.BaseURL() + "/storage/" + fileName

	return fileURL, nil
}

func generateFileName() string {
	h := md5.New()
	randomNumber := fmt.Sprint(rand.Intn(10))
	timestamp := fmt.Sprint(time.Now().Unix())

	h.Write([]byte(randomNumber + timestamp))

	return hex.EncodeToString(h.Sum(nil))
}

func (s *submissionUsecase) FindAllByAssignmentID(ctx context.Context, cursor model.Cursor, assignmentID int64) (submissions []*model.Submission, count int64, err error) {
	submissions, count, err = s.submissionRepo.FindAllByAssignmentID(ctx, cursor, assignmentID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":          utils.Dump(ctx),
			"cursor":       cursor,
			"assignmentID": assignmentID,
		}).Error(err)
		return nil, 0, err
	}

	return
}
