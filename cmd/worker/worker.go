package main

import (
	"os"
	"os/signal"

	"github.com/miun173/autograd/config"
	"github.com/miun173/autograd/worker"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Info("starting worker")
	redisPool := config.NewRedisPool(config.RedisWorkerHost())
	cfg := worker.NewConfig(worker.WithWorkerPool(redisPool))
	wrk := worker.NewWorker(cfg)
	wrk.Start()

	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

	logrus.Info("stopping worker")
	wrk.Stop()
}
