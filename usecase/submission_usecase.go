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
	Upload(ctx context.Context, upload *model.Upload) error
	FindByAssignmentID(ctx context.Context, cursor *model.Cursor, assignmentID int64) ([]*model.Submission, error)
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

func (s *submissionUsecase) Upload(ctx context.Context, upload *model.Upload) error {
	if upload == nil {
		return errors.New("invalid arguments")
	}

	logger := logrus.WithFields(logrus.Fields{
		"ctx":    utils.Dump(ctx),
		"upload": utils.Dump(upload),
	})

	cwd, err := os.Getwd()
	if err != nil {
		logger.Error(err)
		return err
	}

	fileName := generateFileName() + ".cpp"
	filePath := path.Join(cwd, "submission", fileName)
	file, err := os.Create(filePath)
	if err != nil {
		logger.Error(err)
		return err
	}

	file.WriteString(upload.SourceCode)
	defer file.Close()

	upload.FileURL = config.BaseURL() + "/storage/" + fileName

	return nil
}

func generateFileName() string {
	h := md5.New()
	randomNumber := fmt.Sprint(rand.Intn(10))
	timestamp := fmt.Sprint(time.Now().Unix())

	h.Write([]byte(randomNumber + timestamp))

	return hex.EncodeToString(h.Sum(nil))
}

func (s *submissionUsecase) FindByAssignmentID(ctx context.Context, cursor *model.Cursor, assignmentID int64) (submissions []*model.Submission, err error) {
	cursor.Offset = (cursor.Page - 1) * cursor.Size
	submissions, err = s.submissionRepo.FindByAssignmentID(ctx, assignmentID, cursor)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":          utils.Dump(ctx),
			"assignmentID": assignmentID,
		}).Error(err)
		return nil, err
	}

	return
}
