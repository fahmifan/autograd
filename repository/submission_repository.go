package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/miun173/autograd/model"
	"github.com/miun173/autograd/utils"
	"github.com/sirupsen/logrus"
)

// SubmissionRepository ..
type SubmissionRepository interface {
	Create(ctx context.Context, submission *model.Submission) error
	FindByAssignmentID(ctx context.Context, id int64) ([]*model.Submission, error)
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

func (s *submissionRepo) FindByAssignmentID(ctx context.Context, id int64) ([]*model.Submission, error) {
	submissions := []*model.Submission{}
	query := s.db.Where("assignment_id = ?", id).Find(&submissions)
	err := query.Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx": utils.Dump(ctx),
			"id":  id,
		}).Error(err)
		return nil, err
	}

	rows := query.RowsAffected
	if rows == 0 {
		return nil, errors.New("submission with assignmentID " + utils.Int64ToString(id) + " doesn't exist")
	}

	return submissions, nil
}
