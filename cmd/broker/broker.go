package main

import (
	"github.com/fahmifan/autograd/config"
	"github.com/fahmifan/autograd/worker"
)

// for testing wokrer & broker
func main() {
	redisPool := config.NewRedisPool(config.RedisWorkerHost())
	_ = worker.NewBroker(redisPool)
}
