package worker

import (
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

// Config ..
type Config struct {
	pool      *work.WorkerPool
	redisPool *redis.Pool
	enqueuer  *work.Enqueuer
}

// NewConfig ..
func NewConfig(opts ...Option) *Config {
	cfg := &Config{}
	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}

// Option ..
type Option func(*Config)

// WithWorkerPool ..
func WithWorkerPool(rd *redis.Pool) Option {
	return func(c *Config) {
		c.redisPool = rd
	}
}
