package dbconn

import (
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// MustSQLite ..
func MustSQLite() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("autograd.db?_pragma=foreign_keys=true&_journal_mode=WAL"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	return db
}
