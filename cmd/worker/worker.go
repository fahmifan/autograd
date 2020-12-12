package main

import (
	"os"
	"os/signal"

	"github.com/miun173/autograd/config"
	"github.com/miun173/autograd/grader"
	"github.com/miun173/autograd/worker"
	"github.com/sirupsen/logrus"
)

func main() {
	cppCompiler := grader.NewCompiler(grader.CPPCompiler)
	logrus.Info("starting worker")
	redisPool := config.NewRedisPool(config.RedisWorkerHost())
	wrk := worker.NewWorker(
		worker.WithWorkerPool(redisPool),
		worker.WithGrader(grader.NewGrader(cppCompiler)),
	)
	wrk.Start()

	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

	logrus.Info("stopping worker")
	wrk.Stop()
}
