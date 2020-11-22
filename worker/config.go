package worker

import (
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

// config
type cfg struct {
	pool       *work.WorkerPool
	redisPool  *redis.Pool
	enqueuer   *work.Enqueuer
	grader     Grader
	submission Submission
}

// Option ..
type Option func(*cfg)

// WithWorkerPool ..
func WithWorkerPool(rd *redis.Pool) Option {
	return func(c *cfg) {
		c.redisPool = rd
	}
}

// WithGrader ..
func WithGrader(gr Grader) Option {
	return func(c *cfg) {
		c.grader = gr
	}
}
