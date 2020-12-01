package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/miun173/autograd/dto"
	"github.com/miun173/autograd/model"
	"github.com/miun173/autograd/utils"
	"github.com/sirupsen/logrus"
)

// SubmissionRepository ..
type SubmissionRepository interface {
	Create(ctx context.Context, submission *model.Submission) error
	FindByAssignmentID(ctx context.Context, assignmentID int64, pagination *dto.Pagination) ([]*model.Submission, error)
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

func (s *submissionRepo) FindByAssignmentID(ctx context.Context, assignmentID int64, pagination *dto.Pagination) ([]*model.Submission, error) {
	submissions := []*model.Submission{}
	query := s.db.Find(&submissions)
	err := query.Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":          utils.Dump(ctx),
			"assignmentID": assignmentID,
			"pagination":   pagination,
		}).Error(err)
		return nil, err
	}

	rows := query.RowsAffected
	if rows < pagination.Offset {
		return nil, errors.New("page " + utils.Int64ToString(pagination.Page) + " is out of bounds for limit " + utils.Int64ToString(pagination.Limit))
	}

	limit, offset := utils.Int64ToInt(pagination.Limit), utils.Int64ToInt(pagination.Offset)
	query = s.db.Where("assignment_id = ?", assignmentID).Limit(limit).Offset(offset).Order(pagination.Sort).Find(&submissions)
	err = query.Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":          utils.Dump(ctx),
			"assignmentID": assignmentID,
			"pagination":   pagination,
		}).Error(err)
		return nil, err
	}

	rows = query.RowsAffected
	if rows == 0 {
		return nil, errors.New("submission with assignmentID " + utils.Int64ToString(assignmentID) + " doesn't exist")
	}

	return submissions, nil
}
