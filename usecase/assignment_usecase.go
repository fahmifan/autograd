package usecase

import (
	"context"

	"github.com/fahmifan/autograd/model"
	"github.com/fahmifan/autograd/repository"
	"github.com/fahmifan/autograd/utils"
	"github.com/sirupsen/logrus"
)

// AssignmentUsecase ..
type AssignmentUsecase struct {
	assignmentRepo repository.AssignmentRepository
	submissionRepo repository.SubmissionRepository
}

// NewAssignmentUsecase ..
func NewAssignmentUsecase(assignmentRepo repository.AssignmentRepository,
	submissionRepo repository.SubmissionRepository) *AssignmentUsecase {
	return &AssignmentUsecase{
		assignmentRepo: assignmentRepo,
		submissionRepo: submissionRepo,
	}
}

func (a *AssignmentUsecase) Create(ctx context.Context, assignment *model.Assignment) error {
	if assignment == nil {
		return ErrInvalidArguments
	}

	err := a.assignmentRepo.Create(ctx, assignment)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":        utils.Dump(ctx),
			"assignment": utils.Dump(assignment),
		}).Error(err)
		return err
	}

	return nil
}

func (a *AssignmentUsecase) DeleteByID(ctx context.Context, id string) (*model.Assignment, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx": utils.Dump(ctx),
		"id":  id,
	})

	assignment, err := a.FindByID(ctx, id)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	err = a.assignmentRepo.DeleteByID(ctx, id)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return assignment, nil
}

func (a *AssignmentUsecase) FindAll(ctx context.Context, cursor model.Cursor) (assignments []*model.Assignment, count int64, err error) {
	return a.assignmentRepo.FindAll(ctx, cursor)
}

func (a *AssignmentUsecase) FindByID(ctx context.Context, id string) (*model.Assignment, error) {
	assignment, err := a.assignmentRepo.FindByID(ctx, id)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx": utils.Dump(ctx),
			"id":  id,
		}).Error(err)
		return nil, err
	}

	if assignment == nil {
		return nil, ErrNotFound
	}

	return assignment, nil
}

func (a *AssignmentUsecase) FindSubmissionsByID(ctx context.Context, cursor model.Cursor, id string) (submissions []*model.Submission, count int64, err error) {
	submissions, count, err = a.submissionRepo.FindAllByAssignmentID(ctx, cursor, id)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":    utils.Dump(ctx),
			"cursor": utils.Dump(cursor),
			"id":     id,
		}).Error(err)
		return nil, 0, err
	}

	if submissions == nil {
		return nil, 0, ErrNotFound
	}

	return
}

func (a *AssignmentUsecase) Update(ctx context.Context, assignment *model.Assignment) error {
	if assignment == nil {
		return ErrInvalidArguments
	}

	logger := logrus.WithFields(logrus.Fields{
		"ctx":        utils.Dump(ctx),
		"assignment": utils.Dump(assignment),
	})

	_, err := a.FindByID(ctx, assignment.ID)
	if err != nil {
		logger.Error(err)
		return err
	}

	err = a.assignmentRepo.Update(ctx, assignment)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
