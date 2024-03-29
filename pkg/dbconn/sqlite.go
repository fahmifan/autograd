package dbconn

import (
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// MustSQLite ..
func MustSQLite() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("autograd.db?sqlite_foreign_keys=true&_journal_mode=WAL"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.Exec(`PRAGMA journal_mode=WAL;`).Error
	if err != nil {
		log.Fatal(err)
	}

	return db
}
