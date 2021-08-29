package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/fahmifan/autograd/config"
	db "github.com/fahmifan/autograd/db/migrations"
	"github.com/fahmifan/autograd/fs"
	"github.com/fahmifan/autograd/gocrafts"
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

// @title Autograde API
// @version 1.0
// @description API documentation for Autograde

// @BasePath /
func main() {
	redisPool := config.NewRedisPool(config.RedisWorkerHost())
	postgres := db.MustPostgres()
	workerNamespace := config.WorkerNamespace()
	broker := gocrafts.NewBroker(workerNamespace, redisPool)
	localStorer := fs.NewLocalStorer(&fs.Config{
		RootDir: config.FileUploadPath(),
	})

	userRepo := repository.NewUserRepository(postgres)
	userUsecase := usecase.NewUserUsecase(userRepo)
	submissionRepo := repository.NewSubmissionRepo(postgres)
	assignmentRepo := repository.NewAssignmentRepository(postgres)
	sessionRepo := repository.NewSessionRepository(postgres)

	assignmentUsecase := usecase.NewAssignmentUsecase(assignmentRepo, submissionRepo)
	submissionUsecase := usecase.NewSubmissionUsecase(submissionRepo, usecase.SubmissionUsecaseWithBroker(broker))
	graderUsecase := usecase.NewGraderUsecase(submissionUsecase, assignmentUsecase)

	apiServer := httpsvc.NewServer(
		config.Port(),
		config.FileUploadPath(),
		httpsvc.WithUserUsecase(userUsecase),
		httpsvc.WithAssignmentUsecase(assignmentUsecase),
		httpsvc.WithSubmissionUsecase(submissionUsecase),
		httpsvc.WithObjectStorer(localStorer),
		httpsvc.WithSessionRepository(sessionRepo),
	)
	workers := worker.New(&worker.Config{
		Broker:            broker,
		GraderUsecase:     graderUsecase,
		SubmissionUsecase: submissionUsecase,
		AssignmentUsecase: assignmentUsecase,
	})
	workerManager := gocrafts.NewWorkerManager(workerNamespace, config.WorkerConcurrency(), redisPool, workers)

	go func() {
		logrus.Info("run worker")
		workerManager.Start()
	}()

	go func() {
		logrus.Info("run api server")
		apiServer.Run()
	}()

	// gracefull shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	logrus.Info("stopping api server")
	apiServer.Stop(ctx)

	// if worker unable to stop after 1 minute, kill it
	logrus.Info("stopping worker")
	time.AfterFunc(1*time.Minute, func() {
		logrus.Error("unable to stop worker gracefully")
		os.Exit(1)
	})
	workerManager.Stop()
	logrus.Info("worker stopped")
}
