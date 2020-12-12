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
func WithGrader(gr GraderUsecase) Option {
	return func(c *Worker) {
		c.grader = gr
	}
}
