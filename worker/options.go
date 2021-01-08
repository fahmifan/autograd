package worker

import (
	"github.com/gomodule/redigo/redis"
)

// Option ..
type Option func(*Worker)

// WithWorkerPool ..
func WithWorkerPool(rd *redis.Pool) Option {
	return func(c *Worker) {
		c.redisPool = rd
	}
}

// WithGrader ..
func WithGrader(gr Grader) Option {
	return func(c *Worker) {
		c.grader = gr
	}
}

// WithSubmissionUsecase ..
func WithSubmissionUsecase(s SubmissionUsecase) Option {
	return func(c *Worker) {
		c.submission = s
	}
}

// WithAssignmentUsecase ..
func WithAssignmentUsecase(a AssignmentUsecase) Option {
	return func(c *Worker) {
		c.assignment = a
	}
}
