package main

import (
	"os"

	"github.com/miun173/autograd/config"
	db "github.com/miun173/autograd/db/migrations"
	"github.com/miun173/autograd/fs"
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

	exampleRepo := repository.NewExampleRepo()
	exampleUsecase := usecase.NewExampleUsecase(exampleRepo)

	userRepo := repository.NewUserRepository(postgres)
	userUsecase := usecase.NewUserUsecase(userRepo)
	submissionRepo := repository.NewSubmissionRepo(postgres)
	assignmentRepo := repository.NewAssignmentRepository(postgres)
	assignmentUsecase := usecase.NewAssignmentUsecase(assignmentRepo, submissionRepo)
	submissionUsecase := usecase.NewSubmissionUsecase(submissionRepo, usecase.SubmissionUsecaseWithBroker(broker))
	localFS := fs.NewFileSaver()

	server := httpsvc.NewServer(
		config.Port(),
		httpsvc.WithExampleUsecase(exampleUsecase),
		httpsvc.WithUserUsecase(userUsecase),
		httpsvc.WithAssignmentUsecase(assignmentUsecase),
		httpsvc.WithSubmissionUsecase(submissionUsecase),
		httpsvc.WithUploader(localFS),
	)

	server.Run()
}
