package db

import (
	"github.com/fahmifan/autograd/config"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// MustPostgres ..
func MustPostgres() *gorm.DB {
	db, err := gorm.Open(postgres.Open(config.PostgresDSN()), &gorm.Config{})
	if err != nil {
		logrus.Fatal(err)
	}

	return db
}
