package worker

import (
	"github.com/gocraft/work"
	"github.com/miun173/autograd/config"
)

var defaultJobOpt = work.JobOptions{MaxConcurrency: 3, MaxFails: 3}

// Worker :nodoc:
type Worker struct {
	*Config
}

// NewWorker :nodoc:
func NewWorker(cfg *Config) *Worker {
	wrk := &Worker{cfg}
	return wrk
}

// Start starts worker
func (p *Worker) Start() {
	p.registerJobs()
	p.pool.Start()
}

// Stop stops worker
func (p *Worker) Stop() {
	p.pool.Stop()
}

func (p *Worker) registerJobs() {
	conc := config.WorkerConcurrency()
	nameSpace := config.WorkerNamespace()

	p.pool = work.NewWorkerPool(jobHandler{}, conc, nameSpace, p.redisPool)
	p.pool.Middleware(p.registerJobConfig)
	p.pool.JobWithOptions(jobRunCode, defaultJobOpt, (*jobHandler).handleRunCode)
}

func (p *Worker) registerJobConfig(jb *jobHandler, job *work.Job, next work.NextMiddlewareFunc) error {
	jb.Config = p.Config
	return next()
}
