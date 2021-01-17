package worker

import (
	"errors"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/miun173/autograd/config"
	"github.com/miun173/autograd/model"
)

var defaultJobOpt = work.JobOptions{MaxConcurrency: 3, MaxFails: 3}

// The real one
// const cronEvery10Minute = "0 10 * * *"

// TODO: the value is in debug mode
const cronEvery10Minute = "*/60 * * * *"

// Worker ..
type Worker struct {
	pool       *work.WorkerPool
	redisPool  *redis.Pool
	enqueuer   *work.Enqueuer
	grader     model.GraderUsecase
	submission model.SubmissionUsecase
	assignment model.AssignmentUsecase
}

// NewWorker ..
func NewWorker(opts ...Option) *Worker {
	wrk := &Worker{}
	for _, opt := range opts {
		opt(wrk)
	}
	wrk.enqueuer = newEnqueuer(wrk.redisPool)

	return wrk
}

// Start starts worker
func (w *Worker) Start() {
	w.registerJobs()
	w.pool.Start()
}

// Stop stops worker
func (w *Worker) Stop() {
	w.pool.Stop()
}

func (w *Worker) registerJobs() {
	conc := config.WorkerConcurrency()
	nameSpace := config.WorkerNamespace()

	w.pool = work.NewWorkerPool(jobHandler{}, conc, nameSpace, w.redisPool)
	w.pool.Middleware(w.registerJobConfig)

	w.pool.JobWithOptions(jobGradeSubmission, defaultJobOpt, (*jobHandler).handleGradeSubmission)

	// TODO: disable for now
	// w.pool.JobWithOptions(jobGradeAssignment, defaultJobOpt, (*jobHandler).handleGradeAssignment)
	// w.pool.JobWithOptions(jobCheckAllDueAssignments, defaultJobOpt, (*jobHandler).handleCheckAllDueAssignments)
	// w.pool.PeriodicallyEnqueue(cronEvery10Minute, jobCheckAllDueAssignments)
}

func (w *Worker) registerJobConfig(handler *jobHandler, job *work.Job, next work.NextMiddlewareFunc) error {
	if handler == nil {
		return errors.New("unexpected nil handler")
	}

	handler.pool = w.pool
	handler.redisPool = w.redisPool
	handler.enqueuer = w.enqueuer
	handler.grader = w.grader
	handler.submission = w.submission
	handler.assignment = w.assignment

	return next()
}
