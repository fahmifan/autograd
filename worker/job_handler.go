package worker

import (
	"fmt"

	"github.com/gocraft/work"
	"github.com/sirupsen/logrus"
)

// job names
const (
	jobCheckAllDueAssignments string = "check_due_assignment"
	jobGradeAssignment        string = "grade_assignment"
)

// Grader ..
type Grader interface {
	GradeAssignment(assignmentID int64) error
}

type jobHandler struct {
	*cfg
}

func (h *jobHandler) handleCheckAllDueAssignments(job *work.Job) error {
	ids, err := getAllDueAssignments()
	if err != nil {
		logrus.Error(err)
		return fmt.Errorf("unable to get all due assignments: %w", err)
	}

	for _, id := range ids {
		_, err := h.enqueuer.EnqueueUnique(jobGradeAssignment, work.Q{"assignmentID": id})
		if err != nil {
			logrus.Errorf("unable to enqueue assignment %d: %w", id, err)
		}
	}

	return nil
}

func (h *jobHandler) handleGradeAssignment(job *work.Job) error {
	assignmentID := job.ArgInt64("assignmentID")
	h.grader.GradeAssignment(assignmentID)
	return nil
}

func getAllDueAssignments() (ids []int64, err error) {
	return
}
