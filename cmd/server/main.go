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
	"github.com/fahmifan/autograd/web"
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
	postgres := db.MustPostgres()
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

	apiServer := httpsvc.NewServer(
		config.Port(),
		config.FileUploadPath(),
		httpsvc.WithUserUsecase(userUsecase),
		httpsvc.WithAssignmentUsecase(assignmentUsecase),
		httpsvc.WithSubmissionUsecase(submissionUsecase),
		httpsvc.WithMediaUsecase(mediaUsecase),
	)
	wrk := worker.NewWorker(redisPool, worker.WithGrader(graderUsecase))

	debugMode := config.Env() == "development"
	webServer := web.NewServer(config.WebPort(), debugMode)

	go func() {
		logrus.Info("run worker")
		wrk.Start()
	}()

	go func() {
		logrus.Info("run api server")
		apiServer.Run()
	}()

	go func() {
		logrus.Info("run web server")
		webServer.Run()
	}()

	// gracefull shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	logrus.Info("stopping api server")
	apiServer.Stop(ctx)

	logrus.Info("stopping web server")
	webServer.Stop(ctx)

	// if worker unable to stop after 30 seconds kill it
	logrus.Info("stopping worker")
	time.AfterFunc(time.Second*30, func() {
		logrus.Error("unable to stop worker gracefully")
		os.Exit(1)
	})
	wrk.Stop()
	logrus.Info("worker stopped")
}
