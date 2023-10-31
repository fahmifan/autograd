package db

import (
	"github.com/glebarez/sqlite"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// MustSQLite ..
func MustSQLite() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("autograd.db?_pragma=foreign_keys=true&_journal_mode=WAL"), &gorm.Config{})
	if err != nil {
		logrus.Fatal(err)
	}

	return db
}
