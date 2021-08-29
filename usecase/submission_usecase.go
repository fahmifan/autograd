package usecase

import (
	"context"

	"github.com/fahmifan/autograd/model"
	"github.com/fahmifan/autograd/repository"
	"github.com/fahmifan/autograd/utils"

	"github.com/sirupsen/logrus"
)

type SubmissionUsecase struct {
	submissionRepo repository.SubmissionRepository
	workerBroker   model.Broker
}

// SubmissionOption ..
type SubmissionOption func(s *SubmissionUsecase)

// SubmissionUsecaseWithBroker ..
func SubmissionUsecaseWithBroker(b model.Broker) SubmissionOption {
	return func(s *SubmissionUsecase) {
		s.workerBroker = b
	}
}

// NewSubmissionUsecase ..
func NewSubmissionUsecase(submissionRepo repository.SubmissionRepository, opts ...SubmissionOption) *SubmissionUsecase {
	s := &SubmissionUsecase{
		submissionRepo: submissionRepo,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *SubmissionUsecase) Create(ctx context.Context, submission *model.Submission) error {
	if submission == nil {
		return ErrInvalidArguments
	}

	err := s.submissionRepo.Create(ctx, submission)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":        utils.Dump(ctx),
			"submission": utils.Dump(submission),
		}).Error(err)
		return err
	}

	err = s.workerBroker.GradeSubmission(submission.ID)
	if err != nil {
		logrus.Error(err)
	}

	return nil
}

func (s *SubmissionUsecase) DeleteByID(ctx context.Context, id string) (*model.Submission, error) {
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

func (s *SubmissionUsecase) FindByID(ctx context.Context, id string) (*model.Submission, error) {
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

func (s *SubmissionUsecase) Update(ctx context.Context, submission *model.Submission) error {
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

	go func(submissionID string) {
		err := s.workerBroker.GradeSubmission(submissionID)
		if err != nil {
			logger.Error(err)
		}
	}(submission.ID)

	return nil
}

func (s *SubmissionUsecase) FindAllByAssignmentID(ctx context.Context, cursor model.Cursor, assignmentID string) (submissions []*model.Submission, count int64, err error) {
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
func (s *SubmissionUsecase) UpdateGradeByID(ctx context.Context, id string, grade int64) error {
	sbm, err := s.submissionRepo.FindByID(ctx, id)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"id":    id,
			"grade": grade,
		}).Error(err)
		return err
	}

	if sbm == nil {
		return ErrNotFound
	}

	return s.submissionRepo.UpdateGradeByID(ctx, id, grade)
}

// FindByIDAndSubmitter ..
func (s *SubmissionUsecase) FindByIDAndSubmitter(ctx context.Context, id, submitterID string) (*model.Submission, error) {
	sbm, err := s.submissionRepo.FindByIDAndSubmitter(ctx, id, submitterID)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	if sbm == nil {
		return nil, ErrNotFound
	}
	return sbm, nil
}

// FindByAssignmentIDAndSubmitterID ..
func (s *SubmissionUsecase) FindByAssignmentIDAndSubmitterID(ctx context.Context, assignmentID, submitterID string) ([]*model.Submission, error) {
	sbm, err := s.submissionRepo.FindByAssignmentIDAndSubmitterID(ctx, assignmentID, submitterID)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	if sbm == nil {
		return nil, ErrNotFound
	}
	return sbm, nil
}
