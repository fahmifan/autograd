package worker

import (
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/miun173/autograd/config"
)

// newEnqueuer ..
func newEnqueuer(redisPool *redis.Pool) *work.Enqueuer {
	return work.NewEnqueuer(config.WorkerNamespace(), redisPool)
}

// Broker enqueue job for worker
type Broker struct {
	enqueuer *work.Enqueuer
}

// NewBroker ..
func NewBroker(redisPool *redis.Pool) *Broker {
	return &Broker{enqueuer: newEnqueuer(redisPool)}
}
