package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/fahmifan/autograd/config"
	"github.com/fahmifan/autograd/web"
	"github.com/labstack/echo/v4"
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
	ec := echo.New()
	server := web.Server{
		Echo: ec,
		Port: config.WebPort(),
	}

	go func() {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt)
		<-signalChan

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		err := server.Stop(ctx)
		if err != nil {
			logrus.Error(err)
		}
	}()

	logrus.Error(server.Run())
}
