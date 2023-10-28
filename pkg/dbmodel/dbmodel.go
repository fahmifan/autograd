package dbmodel

import (
	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

type Base struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key;"`
	Metadata
}

type Metadata struct {
	CreatedAt null.Time
	UpdatedAt null.Time
	DeletedAt gorm.DeletedAt `sql:"index" json:"deleted_at"`
}

type User struct {
	Base
	Name     string
	Email    string
	Password string
	Role     string
	Active   int
}

type Assignment struct {
	Base
	AssignedBy       uuid.UUID
	Name             string
	Description      string
	CaseInputFileID  uuid.UUID
	CaseOutputFileID uuid.UUID
}

type Submission struct {
	Base
	AssignmentID uuid.UUID
	FileID       uuid.UUID
	SubmittedBy  uuid.UUID
	Grade        int64
	Feedback     string
}

type FileExt string
type FileType string

const (
	FileTypeAssignmentCaseInput  FileType = "assignment_case_input"
	FileTypeAssignmentCaseOutput FileType = "assignment_case_output"
	FileTypeSubmission           FileType = "submission"
)

type File struct {
	Base
	Name string
	Type FileType
	Ext  FileExt
	URL  string
}
