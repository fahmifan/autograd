package main

import (
	"github.com/miun173/autograd/config"
	"github.com/miun173/autograd/worker"
)

// for testing wokrer & broker
func main() {
	redisPool := config.NewRedisPool(config.RedisWorkerHost())
	_ = worker.NewBroker(redisPool)
}
