package worker

import (
	"github.com/gomodule/redigo/redis"
	"github.com/miun173/autograd/model"
	"github.com/miun173/autograd/utils"

	"github.com/gocraft/work"
	"github.com/sirupsen/logrus"
)

// job names
const (
	jobCheckAllDueAssignments string = "check_due_assignment"
	jobGradeAssignment        string = "grade_assignment"
	jobGradeSubmission        string = "grade_submission"
)

type jobHandler struct {
	pool       *work.WorkerPool
	redisPool  *redis.Pool
	enqueuer   *work.Enqueuer
	submission model.SubmissionUsecase
	assignment model.AssignmentUsecase
	grader     model.GraderUsecase
}

func (h *jobHandler) handleGradeSubmission(job *work.Job) error {
	submissionID := utils.StringToInt64(job.ArgString("submissionID"))
	logger := logrus.WithField("submissionID", submissionID)

	err := h.grader.GradeBySubmission(submissionID)
	if err != nil {
		logger.Error(err)
	}

	return err
}

// TODO: enable later
// func (h *jobHandler) handleCheckAllDueAssignments(job *work.Job) error {
// 	logrus.Warn("start >>> ", time.Now())
// 	var size, page int64 = 10, 1
// 	cursor := model.NewCursor(size, page, model.SortCreatedAtDesc)
// 	idsChan := make(chan []int64)

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
// 	eg, ctx := errgroup.WithContext(ctx)
// 	defer cancel()

// 	// produce
// 	eg.Go(func() error {
// 		defer close(idsChan)
// 		for {
// 			ids, _, err := h.assignment.FindAllDueDates(cursor)
// 			if err != nil {
// 				logrus.Error(err)
// 				return fmt.Errorf("unable to get all due assignments: %w", err)
// 			}

// 			if len(ids) == 0 {
// 				break
// 			}

// 			idsChan <- ids

// 			page++
// 			cursor.SetPage(page)
// 		}

// 		return nil
// 	})

// 	// consume
// 	eg.Go(func() error {
// 		select {
// 		case <-ctx.Done():
// 			return ctx.Err()
// 		case ids := <-idsChan:
// 			for _, id := range ids {
// 				_, err := h.enqueuer.EnqueueUnique(jobGradeAssignment, work.Q{"assignmentID": id})
// 				if err != nil {
// 					return fmt.Errorf("unable to enqueue assignment %d: %w", id, err)
// 				}
// 			}
// 		}

// 		return nil
// 	})

// 	err := eg.Wait()
// 	if err != nil && err != context.Canceled {
// 		logrus.Error(err)
// 		return err
// 	}

// 	logrus.Warn("done")
// 	return nil
// }

// TODO: enable later
// func (h *jobHandler) handleGradeAssignment(job *work.Job) error {
// 	assignmentID := job.ArgInt64("assignmentID")
// 	ids, _, err := h.submission.FindAllByAssignmentID(con, assignmentID)
// 	if err != nil {
// 		return err
// 	}

// 	for _, id := range ids {
// 		arg := work.Q{"submissionID": id}
// 		if _, err := h.enqueuer.EnqueueUnique(jobGradeSubmission, arg); err != nil {
// 			logrus.WithField("submissionID", id).Error(err)
// 			return err
// 		}
// 	}

// 	return nil
// }
