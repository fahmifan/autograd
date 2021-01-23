package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/fahmifan/autograd/config"
	db "github.com/fahmifan/autograd/db/migrations"
	"github.com/fahmifan/autograd/fs"
	"github.com/fahmifan/autograd/httpsvc"
	"github.com/fahmifan/autograd/repository"
	"github.com/fahmifan/autograd/usecase"
	"github.com/fahmifan/autograd/worker"
	"github.com/sirupsen/logrus"
)

func initLogger() {
	logLevel := logrus.ErrorLevel

	switch config.Env() {
	case "development", "staging":
		logLevel = logrus.InfoLevel
	}

	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:    true,
		DisableSorting: true,
		DisableColors:  false,
	})

	logrus.SetOutput(os.Stdout)
	logrus.SetReportCaller(true)
	logrus.SetLevel(logLevel)

}

func init() {
	initLogger()
}

func main() {
	redisPool := config.NewRedisPool(config.RedisWorkerHost())
	postgres := db.NewPostgres()
	broker := worker.NewBroker(redisPool)
	localStorage := fs.NewLocalStorage()

	userRepo := repository.NewUserRepository(postgres)
	userUsecase := usecase.NewUserUsecase(userRepo)
	submissionRepo := repository.NewSubmissionRepo(postgres)
	assignmentRepo := repository.NewAssignmentRepository(postgres)

	assignmentUsecase := usecase.NewAssignmentUsecase(assignmentRepo, submissionRepo)
	submissionUsecase := usecase.NewSubmissionUsecase(submissionRepo, usecase.SubmissionUsecaseWithBroker(broker))
	mediaUsecase := usecase.NewMediaUsecase(config.FileUploadPath(), localStorage)
	graderUsecase := usecase.NewGraderUsecase(submissionUsecase, assignmentUsecase)

	server := httpsvc.NewServer(
		config.Port(),
		config.FileUploadPath(),
		httpsvc.WithUserUsecase(userUsecase),
		httpsvc.WithAssignmentUsecase(assignmentUsecase),
		httpsvc.WithSubmissionUsecase(submissionUsecase),
		httpsvc.WithMediaUsecase(mediaUsecase),
	)

	wrk := worker.NewWorker(redisPool, worker.WithGrader(graderUsecase))

	go func() {
		logrus.Info("run server")
		server.Run()
	}()

	go func() {
		logrus.Info("run worker")
		wrk.Start()
	}()

	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

	logrus.Info("stopping server")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	server.Stop(ctx)

	logrus.Info("stopping worker")
	time.AfterFunc(time.Second*30, func() {
		os.Exit(1)
	})
	wrk.Stop()
	logrus.Info("worker stopped")

}
