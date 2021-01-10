package repository

import (
	"context"
	"math"
	"time"

	"gorm.io/gorm"

	"github.com/miun173/autograd/model"
	"github.com/miun173/autograd/utils"
	"github.com/sirupsen/logrus"
)

// SubmissionRepository ..
type SubmissionRepository interface {
	Create(ctx context.Context, submission *model.Submission) error
	DeleteByID(ctx context.Context, id int64) error
	FindAllByAssignmentID(ctx context.Context, cursor model.Cursor, assignmentID int64) ([]*model.Submission, int64, error)
	FindByID(ctx context.Context, id int64) (*model.Submission, error)
	Update(ctx context.Context, submission *model.Submission) error
	UpdateGradeByID(ctx context.Context, id, grade int64) error
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
		return err
	}

	return nil
}

func (s *submissionRepo) DeleteByID(ctx context.Context, id int64) error {
	err := s.db.Where("id = ?", id).Delete(&model.Submission{}).Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx": utils.Dump(ctx),
			"id":  id,
		}).Error(err)
		return err
	}

	return nil
}

func (s *submissionRepo) FindAllByAssignmentID(ctx context.Context, cursor model.Cursor, assignmentID int64) ([]*model.Submission, int64, error) {
	count := int64(0)
	err := s.db.Model(model.Submission{}).Where("assignment_id = ?", assignmentID).Count(&count).Error
	if count == 0 {
		return nil, 0, nil
	}

	submissions := []*model.Submission{}
	err = s.db.Where("assignment_id = ?", assignmentID).Limit(int(cursor.GetSize())).
		Offset(int(cursor.GetOffset())).Order("created_at " + cursor.GetSort()).Find(&submissions).Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":          utils.Dump(ctx),
			"cursor":       utils.Dump(cursor),
			"assignmentID": utils.Dump(assignmentID),
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
	err := s.db.Model(&model.Submission{}).Where("id = ?", submission.ID).Updates(submission).Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":        utils.Dump(ctx),
			"submission": utils.Dump(submission),
		}).Error(err)
	}

	return nil
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

// UpdateGradeByID ..
func (s *submissionRepo) UpdateGradeByID(ctx context.Context, id, grade int64) error {
	err := s.db.
		Model(model.Submission{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"grade":      grade,
			"updated_at": time.Now(),
		}).Error
	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}
