package usecase

import (
	"context"

	"github.com/miun173/autograd/model"
	"github.com/miun173/autograd/repository"
	"github.com/miun173/autograd/utils"
	"github.com/sirupsen/logrus"
)

// AssignmentUsecase ..
type AssignmentUsecase interface {
	Create(ctx context.Context, assignment *model.Assignment) error
	DeleteByID(ctx context.Context, id int64) (*model.Assignment, error)
	FindAll(ctx context.Context, cursor model.Cursor) (assignments []*model.Assignment, count int64, err error)
	FindByID(ctx context.Context, id int64) (*model.Assignment, error)
	FindSubmissionsByID(ctx context.Context, cursor model.Cursor,
		assignmentID int64) (submissions []*model.Submission, count int64, err error)
	Update(ctx context.Context, assignment *model.Assignment) error
}

type assignmentUsecase struct {
	assignmentRepo repository.AssignmentRepository
	submissionRepo repository.SubmissionRepository
}

// NewAssignmentUsecase ..
func NewAssignmentUsecase(assignmentRepo repository.AssignmentRepository,
	submissionRepo repository.SubmissionRepository) AssignmentUsecase {
	return &assignmentUsecase{
		assignmentRepo: assignmentRepo,
		submissionRepo: submissionRepo,
	}
}

func (a *assignmentUsecase) Create(ctx context.Context, assignment *model.Assignment) error {
	if assignment == nil {
		return ErrInvalidArguments
	}

	assignment.ID = utils.GenerateID()
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

func (a *assignmentUsecase) DeleteByID(ctx context.Context, id int64) (*model.Assignment, error) {
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

func (a *assignmentUsecase) FindAll(ctx context.Context, cursor model.Cursor) (assignments []*model.Assignment,
	count int64, err error) {
	assignments, count, err = a.assignmentRepo.FindAll(ctx, cursor)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":    utils.Dump(ctx),
			"cursor": utils.Dump(cursor),
		}).Error(err)
		return nil, 0, err
	}

	return
}

func (a *assignmentUsecase) FindByID(ctx context.Context, id int64) (*model.Assignment, error) {
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

func (a *assignmentUsecase) FindSubmissionsByID(ctx context.Context, cursor model.Cursor,
	id int64) (submissions []*model.Submission, count int64, err error) {
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

func (a *assignmentUsecase) Update(ctx context.Context, assignment *model.Assignment) error {
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
