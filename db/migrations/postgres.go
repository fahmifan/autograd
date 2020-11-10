package db

import (
	"github.com/miun173/autograd/config"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewPostgres :nodoc:
func NewPostgres() *gorm.DB {
	db, err := gorm.Open(postgres.Open(config.PostgresDSN()), &gorm.Config{})
	if err != nil {
		logrus.Fatal(err)
	}

	return db
}
