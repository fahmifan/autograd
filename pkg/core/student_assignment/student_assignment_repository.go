package student_assignment

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/fahmifan/autograd/model"
	"github.com/fahmifan/autograd/pkg/core"
	"github.com/fahmifan/autograd/pkg/dbmodel"
	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

type AssignmentReader struct{}

func (AssignmentReader) FindByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (Assignment, error) {
	assignmentModel := dbmodel.Assignment{}
	err := tx.Where("id = ?", id).First(&assignmentModel).Error
	if err != nil {
		return Assignment{}, fmt.Errorf("find assignment: %w", err)
	}

	return Assignment{
		ID:         assignmentModel.ID,
		DeadlineAt: assignmentModel.DeadlineAt,
	}, nil
}

type StudentAssignmentReader struct{}

type FindAllAssignmentRequest struct {
	core.PaginationRequest
	StudentID uuid.UUID
	From      time.Time
	To        time.Time
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

	err = tx.Scopes(req.PaginateScope).Order("updated_at desc").Find(&assignmentModels).Error
	if err != nil {
		return FindAllAssignmentResponse{}, fmt.Errorf("find assignments: %w", err)
	}

	assignmnetIDs := []uuid.UUID{}
	for _, assignmentModel := range assignmentModels {
		assignmnetIDs = append(assignmnetIDs, assignmentModel.ID)
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

	submissionModels := []dbmodel.Submission{}
	err = tx.Where("assignment_id IN (?) and submitted_by = ?", assignmnetIDs, req.StudentID).
		Find(&submissionModels).
		Error
	if err != nil {
		return FindAllAssignmentResponse{}, fmt.Errorf("find submissions: %w", err)
	}

	return FindAllAssignmentResponse{
		Assignments: toStudentAssignments(assignmentModels, assignerModels, submissionModels),
		Pagination: core.Pagination{
			Page:  req.Page,
			Limit: req.Limit,
			Total: int32(count),
		},
	}, nil
}

type FindByIDOption struct {
	StudentID *uuid.UUID
}

type FindByIDOptionFunc func(*FindByIDOption)

func FindByIDWithStudentID(studentID uuid.UUID) FindByIDOptionFunc {
	return func(fbi *FindByIDOption) {
		fbi.StudentID = &studentID
	}
}

type FindStudentAssignmentByIDRequest struct {
	ID            uuid.UUID
	StudentID     uuid.UUID
	WithStudentID bool
}

func (StudentAssignmentReader) FindByID(ctx context.Context, tx *gorm.DB, req FindStudentAssignmentByIDRequest) (
	StudentAssignment, error,
) {
	assignmentModel := dbmodel.Assignment{}
	err := tx.Where("id = ?", req.ID).First(&assignmentModel).Error
	if err != nil {
		return StudentAssignment{}, fmt.Errorf("find assignment: %w", err)
	}

	assignerModel := dbmodel.User{}
	err = tx.Where("id = ?", assignmentModel.AssignedBy).First(&assignerModel).Error
	if err != nil {
		return StudentAssignment{}, fmt.Errorf("find assigner: %w", err)
	}
	submissionModel := dbmodel.Submission{}
	if req.WithStudentID {
		err = tx.Debug().Where("assignment_id = ? and submitted_by = ?", assignmentModel.ID, req.StudentID).
			First(&submissionModel).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return StudentAssignment{}, fmt.Errorf("find submission: %w", err)
		}
	}

	return toStudentAssignment(assignmentModel, assignerModel, submissionModel), nil
}

type StudentSubmissionWriter struct{}

func (StudentSubmissionWriter) CreateSubmission(ctx context.Context, tx *gorm.DB, submission *StudentSubmission) error {
	submissionModel := dbmodel.Submission{
		Base: dbmodel.Base{
			ID:       submission.ID,
			Metadata: submission.ModelMetadata(),
		},
		AssignmentID: submission.Assignment.ID,
		FileID:       submission.SubmissionFile.ID,
		SubmittedBy:  submission.Student.ID,
		Grade:        submission.Grade,
		Feedback:     submission.Feedback,
	}
	err := tx.Create(submissionModel).Error
	if err != nil {
		return fmt.Errorf("create submission: %w", err)
	}

	return nil
}

func (StudentSubmissionWriter) UpdateSubmission(ctx context.Context, tx *gorm.DB, submission *StudentSubmission) error {
	err := tx.Save(dbmodel.Submission{
		Base: dbmodel.Base{
			ID:       submission.ID,
			Metadata: submission.ModelMetadata(),
		},
		AssignmentID: submission.Assignment.ID,
		FileID:       submission.SubmissionFile.ID,
		SubmittedBy:  submission.Student.ID,
		Grade:        submission.Grade,
		Feedback:     submission.Feedback,
	}).Error
	if err != nil {
		return fmt.Errorf("update submission: %w", err)
	}

	return nil
}

type StudentSubmissionReader struct{}

func (StudentSubmissionReader) FindByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (StudentSubmission, error) {
	submissionModel := dbmodel.Submission{}

	err := tx.Where("id = ?", id).Take(&submissionModel).Error
	if err != nil {
		return StudentSubmission{}, fmt.Errorf("find submission: %w", err)
	}

	studentModel := dbmodel.User{}
	err = tx.Where("id = ?", submissionModel.SubmittedBy).First(&studentModel).Error
	if err != nil {
		return StudentSubmission{}, fmt.Errorf("find student: %w", err)
	}

	assignmentModel := dbmodel.Assignment{}
	err = tx.Where("id = ?", submissionModel.AssignmentID).First(&assignmentModel).Error
	if err != nil {
		return StudentSubmission{}, fmt.Errorf("find assignment: %w", err)
	}

	fileModel := dbmodel.File{}
	err = tx.Where("id = ?", submissionModel.FileID).First(&fileModel).Error
	if err != nil {
		return StudentSubmission{}, fmt.Errorf("find file: %w", err)
	}

	return StudentSubmission{
		ID: submissionModel.ID,
		Student: Student{
			ID:   studentModel.ID,
			Name: studentModel.Name,
		},
		Assignment: Assignment{
			ID:            submissionModel.AssignmentID,
			DeadlineAt:    assignmentModel.DeadlineAt,
			HasAssignment: true,
		},
		SubmissionFile: SubmissionFile{
			ID:                fileModel.ID,
			URL:               fileModel.URL,
			Type:              fileModel.Type,
			TimestampMetadata: core.TimestampMetaFromModel(fileModel.Metadata),
		},
		Grade:             submissionModel.Grade,
		Feedback:          submissionModel.Feedback,
		TimestampMetadata: core.TimestampMetaFromModel(submissionModel.Metadata),
	}, nil
}

func toStudentAssignments(
	assignmentModels []dbmodel.Assignment,
	assigners []dbmodel.User,
	submissions []dbmodel.Submission,
) []StudentAssignment {
	mapAssigner := map[uuid.UUID]dbmodel.User{}
	for _, assigner := range assigners {
		mapAssigner[assigner.ID] = assigner
	}

	assignmentTosubmissionMap := map[uuid.UUID]dbmodel.Submission{}
	for _, submission := range submissions {
		assignmentTosubmissionMap[submission.AssignmentID] = submission
	}

	assignments := make([]StudentAssignment, len(assignmentModels))
	for i, assignmentModel := range assignmentModels {
		submission := assignmentTosubmissionMap[assignmentModel.ID]
		assigner := mapAssigner[assignmentModel.AssignedBy]
		assignments[i] = toStudentAssignment(assignmentModel, assigner, submission)
	}
	return assignments
}

func toStudentAssignment(assignmentModel dbmodel.Assignment, assigner dbmodel.User, submission dbmodel.Submission) StudentAssignment {
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
		Submission: StudentSubmissionForAssignment{
			ID:               submission.ID,
			StudentID:        submission.SubmittedBy,
			Grade:            submission.Grade,
			Feedback:         submission.Feedback,
			SubmissionFileID: submission.FileID,
			UpdatedAt:        submission.UpdatedAt.Time,
			IsGraded:         submission.IsGraded == 1,
		},
		HasSubmission: (submission.ID != uuid.Nil && submission.ID.String() != ""),
	}
}

type SubmissionFileReader struct{}

func (SubmissionFileReader) FindByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (SubmissionFile, error) {
	fileModel := dbmodel.File{}
	err := tx.Where("id = ?", id).First(&fileModel).Error
	if err != nil {
		return SubmissionFile{}, fmt.Errorf("find file: %w", err)
	}

	return SubmissionFile{
		ID:   fileModel.ID,
		URL:  fileModel.URL,
		Type: fileModel.Type,
		TimestampMetadata: core.TimestampMetadata{
			CreatedAt: fileModel.CreatedAt.Time,
			UpdatedAt: fileModel.UpdatedAt.Time,
			DeletedAt: null.Time{NullTime: sql.NullTime(fileModel.DeletedAt)},
		},
	}, nil
}
