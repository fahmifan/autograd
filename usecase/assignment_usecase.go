package usecase

import (
	"context"
	"errors"

	"github.com/miun173/autograd/model"
	"github.com/miun173/autograd/repository"
	"github.com/miun173/autograd/utils"
	"github.com/sirupsen/logrus"
)

// AssignmentUsecase ..
type AssignmentUsecase interface {
	Create(ctx context.Context, assignment *model.Assignment) error
	Delete(ctx context.Context, assignment *model.Assignment) error
	FindAll(ctx context.Context, cursor model.Cursor) (assignments []*model.Assignment, count int64, err error)
	Update(ctx context.Context, assignment *model.Assignment) error
}

type assignmentUsecase struct {
	assignmentRepo repository.AssignmentRepository
}

// NewAssignmentUsecase ..
func NewAssignmentUsecase(assignmentRepo repository.AssignmentRepository) AssignmentUsecase {
	return &assignmentUsecase{
		assignmentRepo: assignmentRepo,
	}
}

func (a *assignmentUsecase) Create(ctx context.Context, assignment *model.Assignment) error {
	if assignment == nil {
		return errors.New("invalid arguments")
	}

	logger := logrus.WithFields(logrus.Fields{
		"ctx":        utils.Dump(ctx),
		"assignment": utils.Dump(assignment),
	})

	assignment.ID = utils.GenerateID()
	err := a.assignmentRepo.Create(ctx, assignment)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (a *assignmentUsecase) Delete(ctx context.Context, assignment *model.Assignment) error {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":        utils.Dump(ctx),
		"assignment": utils.Dump(assignment),
	})

	err := a.assignmentRepo.Delete(ctx, assignment)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (a *assignmentUsecase) FindAll(ctx context.Context, cursor model.Cursor) (assignments []*model.Assignment, count int64, err error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":    utils.Dump(ctx),
		"cursor": utils.Dump(cursor),
	})

	assignments, count, err = a.assignmentRepo.FindAll(ctx, cursor)
	if err != nil {
		logger.Error(err)
		return nil, 0, err
	}

	return
}

func (a *assignmentUsecase) Update(ctx context.Context, assignment *model.Assignment) error {
	if assignment == nil {
		return errors.New("invalid arguments")
	}

	logger := logrus.WithFields(logrus.Fields{
		"ctx":        utils.Dump(ctx),
		"assignment": utils.Dump(assignment),
	})

	err := a.assignmentRepo.Update(ctx, assignment)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
