package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fahmifan/autograd/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		logrus.Fatal(err)
	}

	m, err := migrate.New(
		fmt.Sprintf("file://%s/db/migrations", cwd),
		config.PostgresDSN())
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		log.Fatal(err)
	}
}
