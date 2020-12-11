package repository

import (
	"context"
	"time"

	"github.com/miun173/autograd/model"
	"github.com/miun173/autograd/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// AssignmentRepository ..
type AssignmentRepository interface {
	Create(ctx context.Context, assignment *model.Assignment) error
	Delete(ctx context.Context, assignment *model.Assignment) error
	Update(ctx context.Context, assignment *model.Assignment) error
}

type assignmentRepo struct {
	db *gorm.DB
}

// NewAssignmentRepository ..
func NewAssignmentRepository(db *gorm.DB) AssignmentRepository {
	return &assignmentRepo{
		db: db,
	}
}

func (a *assignmentRepo) Create(ctx context.Context, assignment *model.Assignment) error {
	err := a.db.Create(assignment).Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":        utils.Dump(ctx),
			"assignment": utils.Dump(assignment),
		}).Error(err)
	}

	return err
}

func (a *assignmentRepo) Delete(ctx context.Context, assignment *model.Assignment) error {
	err := a.db.First(assignment).Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":        utils.Dump(ctx),
			"assignment": utils.Dump(assignment),
		}).Error(err)
	}

	deletedAt := time.Now()
	assignment.DeletedAt = &deletedAt

	err = a.db.Save(assignment).Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":        utils.Dump(ctx),
			"assignment": utils.Dump(assignment),
		}).Error(err)
	}

	return err
}

func (a *assignmentRepo) Update(ctx context.Context, assignment *model.Assignment) error {
	err := a.db.Model(&model.Assignment{}).Where("id = ?", assignment.ID).Updates(assignment).Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ctx":        ctx,
			"assignment": assignment,
		}).Error(err)
	}

	return err
}
