package student_assignment

import (
	"errors"
	"time"

	"github.com/fahmifan/autograd/pkg/core"
	"github.com/fahmifan/autograd/pkg/dbmodel"
	"github.com/google/uuid"
)

type Assigner struct {
	ID   uuid.UUID
	Name string
}

type Student struct {
	ID     uuid.UUID
	Name   string
	Active bool
}

type StudentAssignment struct {
	ID            uuid.UUID
	Name          string
	Description   string
	Assigner      Assigner
	DeadlineAt    time.Time
	UpdatedAt     time.Time
	Submission    StudentSubmissionForAssignment
	CodeTemplate  string
	HasSubmission bool
}

type StudentSubmissionForAssignment struct {
	ID               uuid.UUID
	StudentID        uuid.UUID
	Grade            int32
	IsGraded         bool
	Feedback         string
	SubmissionFileID uuid.UUID
	UpdatedAt        time.Time
}

type SubmissionFile struct {
	ID   uuid.UUID
	URL  string
	Type dbmodel.FileType
	core.TimestampMetadata
}

type Assignment struct {
	ID            uuid.UUID
	DeadlineAt    time.Time
	Template      string
	HasAssignment bool
}

type StudentSubmission struct {
	ID             uuid.UUID
	Student        Student
	Assignment     Assignment
	SubmissionFile SubmissionFile
	Grade          int32
	Feedback       string

	core.TimestampMetadata
}

type CreateStudentSubmissionRequest struct {
	NewID          uuid.UUID
	Now            time.Time
	Student        Student
	Assignment     Assignment
	SubmissionFile SubmissionFile
}

func SubmitStudentSubmission(req CreateStudentSubmissionRequest) (StudentSubmission, error) {
	if !req.Student.Active {
		return StudentSubmission{}, errors.New("student must active")
	}

	if req.Assignment.HasAssignment {
		return StudentSubmission{}, errors.New("submission already created")
	}

	if req.Now.After(req.Assignment.DeadlineAt) {
		return StudentSubmission{}, errors.New("assignment deadline has passed")
	}

	if req.SubmissionFile.Type != dbmodel.FileTypeSubmission {
		return StudentSubmission{}, errors.New("invalid submission file")
	}

	return StudentSubmission{
		ID:                req.NewID,
		Student:           req.Student,
		Assignment:        req.Assignment,
		SubmissionFile:    req.SubmissionFile,
		TimestampMetadata: core.NewTimestampMeta(req.Now),
	}, nil
}

type UpdateStudentSubmissionRequest struct {
	Now               time.Time
	NewSubmissionFile SubmissionFile
}

func (studentSub StudentSubmission) Resubmit(req UpdateStudentSubmissionRequest) (StudentSubmission, error) {
	if req.Now.After(studentSub.Assignment.DeadlineAt) {
		return StudentSubmission{}, errors.New("assignment deadline has passed")
	}

	if req.NewSubmissionFile.Type != dbmodel.FileTypeSubmission {
		return StudentSubmission{}, errors.New("invalid submission file")
	}

	studentSub.SubmissionFile = req.NewSubmissionFile
	studentSub.UpdatedAt = req.Now
	return studentSub, nil
}
