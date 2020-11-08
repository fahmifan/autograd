package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/sirupsen/logrus"
)

// NewPostgres :nodoc:
func NewPostgres() *gorm.DB {
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
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatal(err)
	}

	return db
}
