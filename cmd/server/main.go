package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/miun173/autograd/config"
	db "github.com/miun173/autograd/db/migrations"
	"github.com/miun173/autograd/fs"
	"github.com/miun173/autograd/grader"
	"github.com/miun173/autograd/httpsvc"
	"github.com/miun173/autograd/repository"
	"github.com/miun173/autograd/usecase"
	"github.com/miun173/autograd/worker"
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
	exampleRepo := repository.NewExampleRepo()
	exampleUsecase := usecase.NewExampleUsecase(exampleRepo)

	userRepo := repository.NewUserRepository(postgres)
	userUsecase := usecase.NewUserUsecase(userRepo)
	submissionRepo := repository.NewSubmissionRepo(postgres)
	assignmentRepo := repository.NewAssignmentRepository(postgres)
	assignmentUsecase := usecase.NewAssignmentUsecase(assignmentRepo, submissionRepo)
	submissionUsecase := usecase.NewSubmissionUsecase(submissionRepo, usecase.SubmissionUsecaseWithBroker(broker))
	mediaUsecase := usecase.NewMediaUsecase(config.FileUploadPath(), localStorage)

	server := httpsvc.NewServer(
		config.Port(),
		config.FileUploadPath(),
		httpsvc.WithExampleUsecase(exampleUsecase),
		httpsvc.WithUserUsecase(userUsecase),
		httpsvc.WithAssignmentUsecase(assignmentUsecase),
		httpsvc.WithSubmissionUsecase(submissionUsecase),
		httpsvc.WithMediaUsecase(mediaUsecase),
	)

	cppCompiler := grader.NewCompiler(grader.CPPCompiler)
	cppGrader := grader.NewGrader(cppCompiler,
		grader.WithAssignmentUsecase(assignmentUsecase),
		grader.WithSubmissionUsecase(submissionUsecase),
	)
	wrk := worker.NewWorker(
		worker.WithWorkerPool(redisPool),
		worker.WithGrader(cppGrader),
	)

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

	logrus.Info("stopping worker")
	time.AfterFunc(time.Second*10, func() {
		os.Exit(1)
	})
	wrk.Stop()
	logrus.Info("worker stopped")
}
