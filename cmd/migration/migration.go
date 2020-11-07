package main

import (
	"fmt"
	"log"
	"os"

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

	dbname := "autograd"
	port := "5432"
	host := "postgres"
	user := "postgres"
	password := "root"
	dsn := fmt.Sprintf(
		"%s://%s:%s@localhost:%s/%s?sslmode=disable",
		host,
		user,
		password,
		port,
		dbname)

	m, err := migrate.New(
		fmt.Sprintf("file://%s/db/migrations", cwd),
		dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		log.Fatal(err)
	}
}
