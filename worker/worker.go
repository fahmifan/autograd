package worker

import (
	"github.com/gocraft/work"
	"github.com/miun173/autograd/config"
)

var defaultJobOpt = work.JobOptions{MaxConcurrency: 3, MaxFails: 3}

const cronEvery10Minute = "*/10 * * * *"

// Worker ..
type Worker struct {
	*cfg
}

// NewWorker ..
func NewWorker(opts ...Option) *Worker {
	wrkCfg := &cfg{}
	for _, opt := range opts {
		opt(wrkCfg)
	}

	wrk := &Worker{wrkCfg}
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

	w.pool.JobWithOptions(jobGradeAssignment, defaultJobOpt, (*jobHandler).handleGradeAssignment)
	w.pool.JobWithOptions(jobCheckAllDueAssignments, defaultJobOpt, (*jobHandler).handleCheckAllDueAssignments)

	w.pool.PeriodicallyEnqueue(cronEvery10Minute, jobCheckAllDueAssignments)
}

func (w *Worker) registerJobConfig(handler *jobHandler, job *work.Job, next work.NextMiddlewareFunc) error {
	handler.cfg = w.cfg
	return next()
}
