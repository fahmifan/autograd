package worker

import (
	"fmt"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/miun173/autograd/utils"
	"github.com/sirupsen/logrus"
)

// job names
const (
	jobCheckAllDueAssignments string = "check_due_assignment"
	jobGradeAssignment        string = "grade_assignment"
	jobGradeSubmission        string = "grade_submission"
)

// Grader ..
type Grader interface {
	GradeSubmission(submissionID int64) error
}

// Submission ..
type Submission interface {
	FindAllUncheckByAssignmentID(assignmentID int64) (count int64, ids []int64, err error)
}

type jobHandler struct {
	pool       *work.WorkerPool
	redisPool  *redis.Pool
	enqueuer   *work.Enqueuer
	grader     Grader
	submission Submission
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
	_, ids, err := h.submission.FindAllUncheckByAssignmentID(assignmentID)
	if err != nil {
		return err
	}

	for _, id := range ids {
		arg := work.Q{"submissionID": id}
		if _, err := h.enqueuer.EnqueueUnique(jobGradeSubmission, arg); err != nil {
			logrus.Error(err)
			return err
		}
	}

	return nil
}

func (h *jobHandler) handleGradeSubmission(job *work.Job) error {
	submissionID := utils.StringToInt64(job.ArgString("submissionID"))
	return h.grader.GradeSubmission(submissionID)
}

func getAllDueAssignments() (ids []int64, err error) {
	return
}
