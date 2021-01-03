package main

import (
	"os"

	"github.com/miun173/autograd/config"
	db "github.com/miun173/autograd/db/migrations"
	"github.com/miun173/autograd/httpsvc"
	"github.com/miun173/autograd/repository"
	"github.com/miun173/autograd/usecase"
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
	postgres := db.NewPostgres()

	exampleRepo := repository.NewExampleRepo()
	exampleUsecase := usecase.NewExampleUsecase(exampleRepo)

	userRepo := repository.NewUserRepository(postgres)
	userUsecase := usecase.NewUserUsecase(userRepo)
	submissionRepo := repository.NewSubmissionRepo(postgres)
	submissionUsecase := usecase.NewSubmissionUsecase(submissionRepo)
	assignmentRepo := repository.NewAssignmentRepository(postgres)
	assignmentUsecase := usecase.NewAssignmentUsecase(assignmentRepo, submissionRepo)

	server := httpsvc.NewServer(
		config.Port(),
		httpsvc.WithExampleUsecase(exampleUsecase),
		httpsvc.WithUserUsecase(userUsecase),
		httpsvc.WithAssignmentUsecase(assignmentUsecase),
		httpsvc.WithSubmissionUsecase(submissionUsecase),
	)

	server.Run()
}
