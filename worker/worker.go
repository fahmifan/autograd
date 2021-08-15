package worker

import (
	"context"
	"sync"
	"time"

	"github.com/fahmifan/autograd/model"
	"github.com/sirupsen/logrus"
)

// Worker is implementation for async processor
// this should be invoked when
// call this from your worker manager e.g.
// 	- github.com/gocraft/work
// 	- github.com/RichardKnop/machinery
type Worker struct {
	*Config
}

type Config struct {
	_          string // enforce
	Broker     model.Broker
	Grader     model.GraderUsecase
	Submission model.SubmissionUsecase
	Assignment model.AssignmentUsecase
}

func New(cfg *Config) *Worker {
	return &Worker{cfg}
}

// GradeSubmission ..
func (w *Worker) GradeSubmission(submissionID string) error {
	err := w.Grader.GradeBySubmission(submissionID)
	if err != nil {
		logrus.WithField("submissionID", submissionID).Error(err)
	}

	return err
}

// GradeAssignment ..
func (w *Worker) GradeAssignment(assignmentID string) error {
	submissionIDChan := make(chan string)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for id := range submissionIDChan {
			err := w.Broker.GradeSubmission(id)
			if err != nil {
				logrus.WithField("submissionID", id).Error(err)
			}
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	var page, size int64 = 1, 10
	for {
		sorter := model.NewSorter(model.SortCreatedAtAsc.String())
		cursor := model.NewCursor(size, page, sorter)
		subms, _, err := w.Submission.FindAllByAssignmentID(ctx, cursor, assignmentID)
		if err != nil {
			logrus.WithField("assignmentID", assignmentID).Error(err)
			break
		}
		if len(subms) == 0 {
			break
		}

		for _, subm := range subms {
			submissionIDChan <- subm.ID
		}
	}
	close(submissionIDChan)
	wg.Wait()

	return nil
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
