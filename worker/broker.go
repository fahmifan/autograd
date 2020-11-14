package worker

import (
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/miun173/autograd/config"
	"github.com/sirupsen/logrus"
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

// EnqueueJobRunCode ..
func (b *Broker) EnqueueJobRunCode(id int64) (err error) {
	_, err = b.enqueuer.Enqueue(jobRunCode, work.Q{"address": "test@example.com", "subject": "hello world", "customer_id": 4})
	if err != nil {
		logrus.Error(err)
		return err
	}
	return
}
