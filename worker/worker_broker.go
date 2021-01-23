package worker

import (
	"fmt"

	"github.com/fahmifan/autograd/config"
	"github.com/fahmifan/autograd/utils"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
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

// EnqueueJobGradeSubmission ..
func (b *Broker) EnqueueJobGradeSubmission(submissionID int64) error {
	arg := work.Q{"submissionID": utils.Int64ToString(submissionID)}
	_, err := b.enqueuer.EnqueueUnique(jobGradeSubmission, arg)
	if err != nil {
		return fmt.Errorf("failed to enqueue %s: %w", jobGradeSubmission, err)
	}
	return nil
}
