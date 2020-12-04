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
	FindAllByAssignmentID(ctx context.Context, assignmentID int64) ([]*model.Submission, error)
	FindCursorByAssignmentID(ctx context.Context, cursor *model.Cursor, assignmentID int64) ([]*model.Submission, error)
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

func (s *submissionRepo) FindAllByAssignmentID(ctx context.Context, assignmentID int64) ([]*model.Submission, error) {
	submissions := []*model.Submission{}
	err := s.db.Where("assignment_id = ?", assignmentID).Find(&submissions).Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":          utils.Dump(ctx),
			"assignmentID": assignmentID,
		}).Error(err)
		return nil, err
	}

	return submissions, nil
}

func (s *submissionRepo) FindCursorByAssignmentID(ctx context.Context, cursor *model.Cursor, assignmentID int64) ([]*model.Submission, error) {
	submissions := []*model.Submission{}
	query := s.db.Where("assignment_id = ?", assignmentID).Limit(int(cursor.Size)).Offset(int(cursor.Offset)).Order(cursor.Sort).Find(&submissions)
	err := query.Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":          utils.Dump(ctx),
			"cursor":       cursor,
			"assignmentID": assignmentID,
		}).Error(err)
		return nil, err
	}

	rows := query.RowsAffected
	if rows == 0 {
		return nil, errors.New("submission with assignmentID " + utils.Int64ToString(assignmentID) + " doesn't exist")
	}

	return submissions, nil
}
