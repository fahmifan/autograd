package main

import (
	"github.com/gomodule/redigo/redis"
	"github.com/miun173/autograd/config"
	"github.com/miun173/autograd/worker"
)

func main() {
	var redisPool = &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(config.RedisWorkerHost())
		},
	}
	broker := worker.NewBroker(redisPool)
	broker.EnqueueJobRunCode(1)
}
