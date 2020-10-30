package main

import (
	"os"

	"github.com/miun173/autograd/config"
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
	exampleRepo := repository.NewExampleRepo()
	exampleUsecase := usecase.NewExampleUsecase(exampleRepo)

	server := httpsvc.NewServer(
		config.Port(),
		httpsvc.WithExampleUsecase(exampleUsecase),
	)

	server.Run()
}
