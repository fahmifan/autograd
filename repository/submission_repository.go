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
	Delete(ctx context.Context, id int64) (*model.Submission, error)
	FindAllByAssignmentID(ctx context.Context, cursor model.Cursor, assignmentID int64) ([]*model.Submission, int64, error)
	FindByID(ctx context.Context, id int64) (*model.Submission, error)
	Update(ctx context.Context, submission *model.Submission) error
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

func (s *submissionRepo) Delete(ctx context.Context, id int64) (*model.Submission, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx": utils.Dump(ctx),
		"id":  id,
	})

	tx := s.db.Begin()
	submission := &model.Submission{}
	err := tx.Where("id = ?", id).Delete(submission).Error
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	err = tx.Unscoped().Where("id = ?", id).First(submission).Error
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	tx.Commit()

	return submission, nil
}

func (s *submissionRepo) FindAllByAssignmentID(ctx context.Context, cursor model.Cursor, assignmentID int64) ([]*model.Submission, int64, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":          utils.Dump(ctx),
		"assignmentID": assignmentID,
	})

	count := int64(0)
	err := s.db.Model(model.Submission{}).Where("assignment_id = ?", assignmentID).Count(&count).Error
	if err != nil {
		logger.Error(err)
		return nil, count, err
	}

	submissions := []*model.Submission{}
	err = s.db.Where("assignment_id = ?", assignmentID).Limit(int(cursor.GetSize())).
		Offset(int(cursor.GetOffset())).Order("created_at " + cursor.GetSort()).Find(&submissions).Error
	if err != nil {
		logger.WithFields(logrus.Fields{
			"cursor": cursor,
		}).Error(err)
		return nil, count, err
	}

	return submissions, count, nil
}

func (s *submissionRepo) FindByID(ctx context.Context, id int64) (*model.Submission, error) {
	submission := &model.Submission{}
	err := s.db.Where("id = ?", id).Take(submission).Error
	switch err {
	case nil: // ignore
	case gorm.ErrRecordNotFound:
		return nil, nil
	default:
		logrus.WithFields(logrus.Fields{
			"ctx": utils.Dump(ctx),
			"id":  id,
		}).Error(err)
		return nil, err
	}

	return submission, nil
}

func (s *submissionRepo) Update(ctx context.Context, submission *model.Submission) error {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":        ctx,
		"submission": submission,
	})

	tx := s.db.Begin()
	err := s.db.Model(&model.Submission{}).Where("id = ?", submission.ID).Updates(submission).Error
	if err != nil {
		logger.Error(err)
		return err
	}

	err = s.db.Where("id = ?", submission.ID).First(submission).Error
	if err != nil {
		logger.Error(err)
		return err
	}

	tx.Commit()

	return err
}
