package student_assignment

import (
	"time"

	"github.com/google/uuid"
)

type StudentAssignment struct {
	ID          uuid.UUID
	Name        string
	Description string
	Assigner    Assigner
	DeadlineAt  time.Time
	UpdatedAt   time.Time
}

type Assigner struct {
	ID   uuid.UUID
	Name string
}
