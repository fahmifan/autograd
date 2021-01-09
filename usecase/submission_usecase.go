package usecase

import (
	"context"
	"crypto/md5"
	"encoding/hex"
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
	DeleteByID(ctx context.Context, id int64) (*model.Submission, error)
	FindByID(ctx context.Context, id int64) (*model.Submission, error)
	Update(ctx context.Context, submission *model.Submission) error
	Upload(ctx context.Context, sourceCode string) (string, error)
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
		return ErrInvalidArguments
	}

	submission.ID = utils.GenerateID()
	err := s.submissionRepo.Create(ctx, submission)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":        utils.Dump(ctx),
			"submission": utils.Dump(submission),
		}).Error(err)
		return err
	}

	return nil
}

func (s *submissionUsecase) DeleteByID(ctx context.Context, id int64) (*model.Submission, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx": utils.Dump(ctx),
		"id":  utils.Dump(id),
	})

	submission, err := s.FindByID(ctx, id)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	err = s.submissionRepo.DeleteByID(ctx, id)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return submission, nil
}

func (s *submissionUsecase) FindByID(ctx context.Context, id int64) (*model.Submission, error) {
	submission, err := s.submissionRepo.FindByID(ctx, id)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx": utils.Dump(ctx),
			"id":  utils.Dump(id),
		}).Error(err)
		return nil, err
	}

	if submission == nil {
		return nil, ErrNotFound
	}

	return submission, nil
}

func (s *submissionUsecase) Update(ctx context.Context, submission *model.Submission) error {
	if submission == nil {
		return ErrInvalidArguments
	}

	err := s.submissionRepo.Update(ctx, submission)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":        utils.Dump(ctx),
			"submission": utils.Dump(submission),
		}).Error(err)
		return err
	}

	return nil
}

func (s *submissionUsecase) Upload(ctx context.Context, sourceCode string) (string, error) {
	if sourceCode == "" {
		return "", ErrInvalidArguments
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
