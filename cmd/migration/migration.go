package main

import (
	"fmt"
	"log"
	"os"

	db "github.com/fahmifan/autograd/db/migrations"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbgorm := db.MustSQLite()
	dbconn, err := dbgorm.DB()
	if err != nil {
		log.Fatal(err)
	}

	driver, err := sqlite3.WithInstance(dbconn, &sqlite3.Config{
		DatabaseName: "autograd.db",
	})
	if err != nil {
		log.Fatal(err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	src := fmt.Sprintf("file://%s/db/migrations", cwd)

	m, err := migrate.NewWithDatabaseInstance(
		src,
		"autograd",
		driver,
	)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		log.Fatal(err)
	}
}
