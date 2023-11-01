package student_assignment

import (
	"context"
	"fmt"
	"time"

	"github.com/fahmifan/autograd/model"
	"github.com/fahmifan/autograd/pkg/core"
	"github.com/fahmifan/autograd/pkg/dbmodel"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StudentAssignmentReader struct{}

type FindAllAssignmentRequest struct {
	core.PaginationRequest
	From time.Time
	To   time.Time
}

func (req FindAllAssignmentRequest) GetFrom(now time.Time) time.Time {
	// show last 7 days assignment by default
	if req.From.IsZero() {
		return now.AddDate(0, 0, -7)
	}

	return req.From
}

func (req FindAllAssignmentRequest) GetTo(now time.Time) time.Time {
	// show next 7 days assignment by default
	if req.To.IsZero() {
		return now.AddDate(0, 0, 7)
	}

	return req.To
}

type FindAllAssignmentResponse struct {
	core.Pagination
	Assignments []StudentAssignment
}

func (StudentAssignmentReader) FindAllAssignments(ctx context.Context, tx *gorm.DB, req FindAllAssignmentRequest) (
	FindAllAssignmentResponse, error,
) {
	assignmentModels := []dbmodel.Assignment{}
	count := int64(0)

	err := tx.Model(&model.Assignment{}).Count(&count).Error
	if err != nil {
		return FindAllAssignmentResponse{}, fmt.Errorf("count assignments: %w", err)
	}

	err = tx.Scopes(req.PaginateScope).Find(&assignmentModels).Error
	if err != nil {
		return FindAllAssignmentResponse{}, fmt.Errorf("find assignments: %w", err)
	}

	assignerIDs := []uuid.UUID{}
	for _, assignmentModel := range assignmentModels {
		assignerIDs = append(assignerIDs, assignmentModel.AssignedBy)
	}

	assignerModels := []dbmodel.User{}
	err = tx.Where("id IN (?) ", assignerIDs).Find(&assignerModels).Error
	if err != nil {
		return FindAllAssignmentResponse{}, fmt.Errorf("find assigners: %w", err)
	}

	return FindAllAssignmentResponse{
		Assignments: toStudentAssignments(assignmentModels, assignerModels),
		Pagination: core.Pagination{
			Page:  req.Page,
			Limit: req.Limit,
			Total: int32(count),
		},
	}, nil
}

func (StudentAssignmentReader) FindByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (
	StudentAssignment, error,
) {
	assignmentModel := dbmodel.Assignment{}
	err := tx.Where("id = ?", id).First(&assignmentModel).Error
	if err != nil {
		return StudentAssignment{}, fmt.Errorf("find assignment: %w", err)
	}

	assignerModel := dbmodel.User{}
	err = tx.Where("id = ?", assignmentModel.AssignedBy).First(&assignerModel).Error
	if err != nil {
		return StudentAssignment{}, fmt.Errorf("find assigner: %w", err)
	}

	return toStudentAssignment(assignmentModel, assignerModel), nil
}

func toStudentAssignments(assignmentModels []dbmodel.Assignment, assigners []dbmodel.User) []StudentAssignment {
	mapAssigner := map[uuid.UUID]dbmodel.User{}
	for _, assigner := range assigners {
		mapAssigner[assigner.ID] = assigner
	}

	assignments := make([]StudentAssignment, len(assignmentModels))
	for i, assignmentModel := range assignmentModels {
		assignments[i] = toStudentAssignment(assignmentModel, mapAssigner[assignmentModel.AssignedBy])
	}
	return assignments
}

func toStudentAssignment(assignmentModel dbmodel.Assignment, assigner dbmodel.User) StudentAssignment {
	return StudentAssignment{
		ID:          assignmentModel.ID,
		Name:        assignmentModel.Name,
		Description: assignmentModel.Description,
		Assigner: Assigner{
			ID:   assignmentModel.AssignedBy,
			Name: assigner.Name,
		},
		DeadlineAt: assignmentModel.DeadlineAt,
		UpdatedAt:  assignmentModel.UpdatedAt.Time,
	}
}
