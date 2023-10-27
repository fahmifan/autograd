package dbmodel

import (
	"time"

	"github.com/fahmifan/autograd/pkg/core/auth"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `sql:"index" json:"deleted_at"`
}

type User struct {
	Base
	Name     string
	Email    string
	Password string
	Role     auth.Role
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
