package worker

import (
	"context"
	"sync"
	"time"

	"github.com/fahmifan/autograd/model"
	"github.com/sirupsen/logrus"
)

// Worker ..
type Worker struct {
	*Config
}

type Config struct {
	_                 string // enforce
	Broker            model.Broker
	GraderUsecase     model.GraderUsecase
	SubmissionUsecase model.SubmissionUsecase
	AssignmentUsecase model.AssignmentUsecase
}

func New(cfg *Config) *Worker {
	return &Worker{cfg}
}

// GradeSubmission ..
func (w *Worker) GradeSubmission(submissionID string) error {
	err := w.GraderUsecase.GradeBySubmission(submissionID)
	if err != nil {
		logrus.WithField("submissionID", submissionID).Error(err)
	}

	return err
}

// GradeAssignment ..
func (w *Worker) GradeAssignment(assignmentID string) error {
	const maxBuf = 10
	submsChan := make(chan []*model.Submission, maxBuf)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go w.gradeSubmissions(wg, submsChan)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	var page int64 = 1
	const size int64 = 10
	for {
		sorter := model.NewSorter(model.SortCreatedAtAsc.String())
		cursor := model.NewCursor(size, page, sorter)
		subms, _, err := w.SubmissionUsecase.FindAllByAssignmentID(ctx, cursor, assignmentID)
		if err != nil {
			logrus.WithField("assignmentID", assignmentID).Error(err)
			return err
		}
		if len(subms) == 0 {
			break
		}

		submsChan <- subms
		page++
	}
	close(submsChan)
	wg.Wait()

	return nil
}

func (w *Worker) gradeSubmissions(wg *sync.WaitGroup, submsChan chan []*model.Submission) {
	defer wg.Done()
	for subms := range submsChan {
		for _, subm := range subms {
			if subm == nil {
				continue
			}
			err := w.Broker.GradeSubmission(subm.ID)
			if err != nil {
				logrus.WithField("submissionID", subm).Error(err)
			}
		}
	}
}
