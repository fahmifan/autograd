package assignments

import (
	"errors"
	"strings"
	"time"

	"github.com/fahmifan/autograd/pkg/core"
	"github.com/fahmifan/autograd/pkg/dbmodel"
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
	core.EntityMeta
}

type Assignment struct {
	ID             uuid.UUID
	Name           string
	Description    string
	Assigner       Assigner
	CaseInputFile  CaseFile
	CaseOutputFile CaseFile

	core.EntityMeta
}

type CreateAssignmentRequest struct {
	NewID          uuid.UUID
	Now            time.Time
	Assigner       Assigner
	Name           string
	Description    string
	CaseInputFile  CaseFile
	CaseOutputFile CaseFile
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

	return Assignment{
		EntityMeta:     core.NewEntityMeta(req.Now),
		Name:           req.Name,
		Description:    req.Description,
		CaseInputFile:  req.CaseInputFile,
		CaseOutputFile: req.CaseOutputFile,
		Assigner:       req.Assigner,
	}, nil
}

type UpdateAssignmentRequest struct {
	Now            time.Time
	Assigner       Assigner
	Name           string
	Description    string
	CaseInputFile  CaseFile
	CaseOutputFile CaseFile
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
	Submitter  Assigner
	SourceFile SubmissionFile
	Grade      int64
	Feedback   string
}
