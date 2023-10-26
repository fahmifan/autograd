package db

import (
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// MustSQLite ..
func MustSQLite() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("autograd.db"), &gorm.Config{})
	if err != nil {
		logrus.Fatal(err)
	}

	return db
}
