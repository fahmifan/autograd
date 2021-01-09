package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/miun173/autograd/config"
	db "github.com/miun173/autograd/db/migrations"
	"github.com/miun173/autograd/grader"
	"github.com/miun173/autograd/model"
	"github.com/miun173/autograd/repository"
	"github.com/miun173/autograd/usecase"
	"github.com/miun173/autograd/worker"
	"github.com/sirupsen/logrus"
)

type assignment struct {
}

var counter int64 = 0

func (a *assignment) FindAllDueDates(c model.Cursor) ([]int64, int64, error) {
	if counter == 3 {
		return nil, 0, nil
	}
	ids := []int64{counter + 1, counter + 2, counter + 3}
	counter++
	return ids, int64(len(ids)), nil
}

func main() {
	cppCompiler := grader.NewCompiler(grader.CPPCompiler)
	logrus.Info("starting worker")
	redisPool := config.NewRedisPool(config.RedisWorkerHost())

	postgres := db.NewPostgres()

	submissionRepo := repository.NewSubmissionRepo(postgres)
	assignmentRepo := repository.NewAssignmentRepository(postgres)

	sumissionUsecase := usecase.NewSubmissionUsecase(submissionRepo)
	assignmentUsecase := usecase.NewAssignmentUsecase(assignmentRepo, submissionRepo)

	cppGrader := grader.NewGrader(cppCompiler,
		grader.WithAssignmentUsecase(assignmentUsecase),
		grader.WithSubmissionUsecase(sumissionUsecase),
	)
	wrk := worker.NewWorker(
		worker.WithWorkerPool(redisPool),
		worker.WithGrader(cppGrader),
	)
	wrk.Start()

	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

	logrus.Info("stopping worker")
	time.AfterFunc(time.Second*10, func() {
		os.Exit(1)
	})
	wrk.Stop()
	logrus.Info("worker stopped")
}
