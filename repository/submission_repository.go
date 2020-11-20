package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/miun173/autograd/model"
	"github.com/miun173/autograd/utils"
	"github.com/sirupsen/logrus"
)

// SubmissionRepository ..
type SubmissionRepository interface {
	Create(ctx context.Context, submission *model.Submission) error
}

type submissionRepo struct {
	db *gorm.DB
}

// NewSubmissionRepo ..
func NewSubmissionRepo(db *gorm.DB) SubmissionRepository {
	return &submissionRepo{
		db: db,
	}
}

func (s *submissionRepo) Create(ctx context.Context, submission *model.Submission) error {
	err := s.db.Create(submission).Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":        utils.Dump(ctx),
			"submission": utils.Dump(submission),
		}).Error(err)
	}

	return err
}
