package assignments

import (
	"errors"
	"strings"
	"time"

	"github.com/fahmifan/autograd/pkg/core"
	"github.com/fahmifan/autograd/pkg/dbmodel"
	"github.com/fahmifan/ulids"
	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

type Assigner struct {
	ID     uuid.UUID
	Name   string
	Active bool
}

type CaseFile struct {
	ID   uuid.UUID
	URL  string
	Type dbmodel.FileType
	core.TimestampMetadata
}

type Course struct {
	ID          ulids.ULID
	Name        string
	Description string
	IsActive    bool
}

type Assignment struct {
	ID             uuid.UUID
	CourseID       ulids.ULID
	Name           string
	Description    string
	Template       string
	Assigner       Assigner
	Course         Course
	CaseInputFile  CaseFile
	CaseOutputFile CaseFile
	DeadlineAt     time.Time

	core.TimestampMetadata
}

type CreateAssignmentRequest struct {
	NewID          uuid.UUID
	Now            time.Time
	Assigner       Assigner
	Course         Course
	Name           string
	Description    string
	Template       string
	CaseInputFile  CaseFile
	CaseOutputFile CaseFile
	DeadlineAt     time.Time
}

func CreateAssignment(req CreateAssignmentRequest) (Assignment, error) {
	if !req.Assigner.Active {
		return Assignment{}, errors.New("assigner must active")
	}

	if len(strings.TrimSpace(req.Name)) < 3 {
		return Assignment{}, errors.New("name must be at least 3 characters")
	}

	if len(strings.TrimSpace(req.Description)) < 3 {
		return Assignment{}, errors.New("description must be at least 3 characters")
	}

	if !req.Course.IsActive {
		return Assignment{}, errors.New("course is not active")
	}

	return Assignment{
		ID:                req.NewID,
		TimestampMetadata: core.NewTimestampMeta(req.Now),
		Name:              req.Name,
		Description:       req.Description,
		CaseInputFile:     req.CaseInputFile,
		CaseOutputFile:    req.CaseOutputFile,
		Assigner:          req.Assigner,
		DeadlineAt:        req.DeadlineAt,
		Template:          req.Template,
		CourseID:          req.Course.ID,
		Course:            req.Course,
	}, nil
}

type UpdateAssignmentRequest struct {
	Now            time.Time
	Assigner       Assigner
	Name           string
	Description    string
	Template       string
	CaseInputFile  CaseFile
	CaseOutputFile CaseFile
	DeadlineAt     time.Time
}

func (assignment Assignment) Update(req UpdateAssignmentRequest) (Assignment, error) {
	if !req.Assigner.Active {
		return Assignment{}, errors.New("assigner must active")
	}

	if len(strings.TrimSpace(req.Name)) < 3 {
		return Assignment{}, errors.New("name must be at least 3 characters")
	}

	if len(strings.TrimSpace(req.Description)) < 3 {
		return Assignment{}, errors.New("description must be at least 3 characters")
	}

	assignment.Name = req.Name
	assignment.Description = req.Description
	assignment.CaseInputFile = req.CaseInputFile
	assignment.CaseOutputFile = req.CaseOutputFile
	assignment.Assigner = req.Assigner
	assignment.UpdatedAt = req.Now
	assignment.DeadlineAt = req.DeadlineAt
	assignment.Template = req.Template

	return assignment, nil
}

func (assignment Assignment) Delete(now time.Time) (Assignment, error) {
	assignment.DeletedAt = null.TimeFrom(now)
	return assignment, nil
}

type SubmissionFile struct {
	ID  uuid.UUID
	URL string
}

type Submission struct {
	ID         uuid.UUID
	Assignment Assignment
	Submitter  Submitter
	SourceFile SubmissionFile
	Grade      int32
	Feedback   string
	core.TimestampMetadata
}

type Submitter struct {
	ID     uuid.UUID
	Name   string
	Active bool
}

type CreateSubmissionRequest struct {
	NewID          uuid.UUID
	Now            time.Time
	Assignment     Assignment
	Submitter      Submitter
	SubmissionFile SubmissionFile
}

func CreateSubmission(req CreateSubmissionRequest) (Submission, error) {
	if req.Now.After(req.Assignment.DeadlineAt) {
		return Submission{}, errors.New("assignment is already closed")
	}

	if !req.Submitter.Active {
		return Submission{}, errors.New("submitter must active")
	}

	subm := Submission{
		ID:                req.NewID,
		Assignment:        req.Assignment,
		Submitter:         req.Submitter,
		SourceFile:        req.SubmissionFile,
		Grade:             0,
		Feedback:          "",
		TimestampMetadata: core.NewTimestampMeta(req.Now),
	}

	return subm, nil
}

type UpdateSubmissionRequest struct {
	Now            time.Time
	Submitter      Submitter
	SubmissionFile SubmissionFile
}

func (submission Submission) Update(req UpdateSubmissionRequest) (Submission, error) {
	if req.Now.After(submission.Assignment.DeadlineAt) {
		return Submission{}, errors.New("assignment is already closed")
	}

	submission.SourceFile = req.SubmissionFile
	submission.UpdatedAt = req.Now
	submission.Submitter = req.Submitter

	return submission, nil
}

func (submission Submission) IsOwner(userID uuid.UUID) bool {
	return submission.Submitter.ID == userID
}

func (submission Submission) Delete(now time.Time) (Submission, error) {
	submission.DeletedAt = null.TimeFrom(now)
	return submission, nil
}
