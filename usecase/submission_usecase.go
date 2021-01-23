package usecase

import (
	"context"

	"github.com/fahmifan/autograd/model"
	"github.com/fahmifan/autograd/repository"
	"github.com/fahmifan/autograd/utils"

	"github.com/sirupsen/logrus"
)

// SubmissionUsecase ..
type SubmissionUsecase interface {
	Create(ctx context.Context, submission *model.Submission) error
	DeleteByID(ctx context.Context, id int64) (*model.Submission, error)
	FindByID(ctx context.Context, id int64) (*model.Submission, error)
	Update(ctx context.Context, submission *model.Submission) error
	FindAllByAssignmentID(ctx context.Context, cursor model.Cursor, assignmentID int64) (submissions []*model.Submission, count int64, err error)
	UpdateGradeByID(ctx context.Context, id, grade int64) error
}

type submissionUsecase struct {
	submissionRepo repository.SubmissionRepository
	assignmentRepo repository.AssignmentRepository
	broker         model.WorkerBroker
}

// SubmissionOption ..
type SubmissionOption func(s *submissionUsecase)

// SubmissionUsecaseWithBroker ..
func SubmissionUsecaseWithBroker(b model.WorkerBroker) SubmissionOption {
	return func(s *submissionUsecase) {
		s.broker = b
	}
}

// NewSubmissionUsecase ..
func NewSubmissionUsecase(submissionRepo repository.SubmissionRepository, opts ...SubmissionOption) SubmissionUsecase {
	s := &submissionUsecase{
		submissionRepo: submissionRepo,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
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

	go func(sbmID int64) {
		err := s.broker.EnqueueJobGradeSubmission(sbmID)
		if err != nil {
			logrus.Error(err)
		}
	}(submission.ID)

	return nil
}

func (s *submissionUsecase) DeleteByID(ctx context.Context, id int64) (*model.Submission, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx": utils.Dump(ctx),
		"id":  id,
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
			"id":  id,
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

	logger := logrus.WithFields(logrus.Fields{
		"ctx":        utils.Dump(ctx),
		"submission": utils.Dump(submission),
	})

	_, err := s.FindByID(ctx, submission.ID)
	if err != nil {
		logger.Error(err)
		return err
	}

	err = s.submissionRepo.Update(ctx, submission)
	if err != nil {
		logger.Error(err)
		return err
	}

	go func(submissionID int64) {
		err := s.broker.EnqueueJobGradeSubmission(submissionID)
		if err != nil {
			logger.Error(err)
		}
	}(submission.ID)

	return nil
}

func (s *submissionUsecase) FindAllByAssignmentID(ctx context.Context, cursor model.Cursor, assignmentID int64) (submissions []*model.Submission, count int64, err error) {
	submissions, count, err = s.submissionRepo.FindAllByAssignmentID(ctx, cursor, assignmentID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":          utils.Dump(ctx),
			"cursor":       utils.Dump(cursor),
			"assignmentID": assignmentID,
		}).Error(err)
		return nil, 0, err
	}

	return
}

// UpdateGradeByID ..
func (s *submissionUsecase) UpdateGradeByID(ctx context.Context, id, grade int64) error {
	sbm, err := s.submissionRepo.FindByID(ctx, id)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"id":    id,
			"grade": grade,
		}).Error(err)
	}

	if sbm == nil {
		return ErrNotFound
	}

	return s.submissionRepo.UpdateGradeByID(ctx, id, grade)
}
