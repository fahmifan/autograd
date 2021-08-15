package gocrafts

import (
	"errors"
	"fmt"

	"github.com/fahmifan/autograd/config"
	"github.com/fahmifan/autograd/model"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

// job names
const (
	jobCheckAllDueAssignments string = "check_all_due_assignment"
	jobGradeAssignment        string = "grade_assignment"
	jobGradeSubmission        string = "grade_submission"
)

// Broker enqueue job for a worker
type Broker struct {
	enqueuer *work.Enqueuer
}

// NewBroker ..
func NewBroker(worksapce string, redisPool *redis.Pool) *Broker {
	return &Broker{enqueuer: work.NewEnqueuer(worksapce, redisPool)}
}

// GradeSubmission ..
func (b *Broker) GradeSubmission(submissionID string) error {
	arg := work.Q{"submissionID": submissionID}
	_, err := b.enqueuer.Enqueue(jobGradeSubmission, arg)
	if err != nil {
		return fmt.Errorf("failed to enqueue %s: %w", jobGradeSubmission, err)
	}
	return nil
}

// GradeAssignment ..
func (b *Broker) GradeAssignment(assignmentID string) error {
	arg := work.Q{"assignmentID": assignmentID}
	_, err := b.enqueuer.Enqueue(jobGradeAssignment, arg)
	if err != nil {
		return fmt.Errorf("failed to enqueue %s: %w", jobGradeSubmission, err)
	}
	return nil
}

// WorkerManager subscribe to broker and assign job to a Worker through workerAdapter
type WorkerManager struct {
	pool       *work.WorkerPool
	redisPool  *redis.Pool
	enqueuer   *work.Enqueuer
	worker     model.Worker
	concurency uint
}

const DefaultConcurrency uint = 3

// NewWorkerManager ..
func NewWorkerManager(namespace string, concurency uint, redisPool *redis.Pool, worker model.Broker) *WorkerManager {
	if concurency <= 0 {
		concurency = DefaultConcurrency
	}
	wrk := &WorkerManager{
		redisPool:  redisPool,
		enqueuer:   work.NewEnqueuer(namespace, redisPool),
		worker:     worker,
		concurency: concurency,
	}

	return wrk
}

// Start starts worker manager
func (rw *WorkerManager) Start() {
	rw.registerJobs()
	rw.pool.Start()
}

// Stop stops worker manager
func (w *WorkerManager) Stop() {
	w.pool.Stop()
}

func (rw *WorkerManager) registerJobs() {
	nameSpace := config.WorkerNamespace()
	defaultJobOpt := work.JobOptions{MaxConcurrency: 3, MaxFails: 3}

	rw.pool = work.NewWorkerPool(workerAdapter{}, rw.concurency, nameSpace, rw.redisPool)
	rw.pool.Middleware(rw.initWorkerAdapter)

	rw.pool.JobWithOptions(jobGradeSubmission, defaultJobOpt, (*workerAdapter).GradeSubmission)
	rw.pool.JobWithOptions(jobGradeAssignment, defaultJobOpt, (*workerAdapter).GradeAssignment)
}

func (rw *WorkerManager) initWorkerAdapter(wrk *workerAdapter, job *work.Job, next work.NextMiddlewareFunc) error {
	if wrk == nil {
		return errors.New("unexpected nil handler")
	}

	wrk.pool = rw.pool
	wrk.redisPool = rw.redisPool
	wrk.worker = rw.worker

	return next()
}

// adapter that call the actual Worker implementation
type workerAdapter struct {
	pool      *work.WorkerPool
	redisPool *redis.Pool
	worker    model.Worker
}

func (h *workerAdapter) GradeSubmission(job *work.Job) error {
	return h.worker.GradeSubmission(job.ArgString("submissionID"))
}

func (h *workerAdapter) GradeAssignment(job *work.Job) error {
	return h.worker.GradeAssignment(job.ArgString("assingmentID"))
}
