package main

import (
	"github.com/fahmifan/autograd/config"
	"github.com/fahmifan/autograd/gocrafts"
)

// for testing worker & broker
func main() {
	redisPool := config.NewRedisPool(config.RedisWorkerHost())
	_ = gocrafts.NewBroker(config.WorkerNamespace(), redisPool)
}
