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
