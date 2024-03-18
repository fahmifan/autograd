package dbmodel

import (
	"time"

	"github.com/fahmifan/ulids"
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

func (user User) IsActive() bool {
	return user.Active == 1
}

type RelUserToActivationToken struct {
	UserID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	ActivationTokenID uuid.UUID `gorm:"type:uuid;primaryKey"`
	DeletedAt         gorm.DeletedAt
}

func (RelUserToActivationToken) TableName() string {
	return "rel_user_to_activation_tokens"
}

type ActivationToken struct {
	Base
	Token     string
	ExpiredAt time.Time
}

func (ActivationToken) TableName() string {
	return "activation_tokens"
}

type Assignment struct {
	Base
	AssignedBy       uuid.UUID
	Name             string
	Description      string
	Template         string
	CaseInputFileID  uuid.UUID
	CaseOutputFileID uuid.UUID
	DeadlineAt       time.Time
}

type Submission struct {
	Base
	AssignmentID uuid.UUID
	FileID       uuid.UUID
	SubmittedBy  uuid.UUID
	Grade        int32
	Feedback     string
	IsGraded     int
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

type OutboxItem struct {
	ID            ulids.ULID
	IdempotentKey string
	Status        string
	JobType       string
	Payload       string
	Version       int32
}
