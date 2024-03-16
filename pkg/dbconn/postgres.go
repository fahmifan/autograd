package dbconn

import (
	"database/sql"

	"github.com/fahmifan/autograd/pkg/xsqlc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func MustPostgres() *gorm.DB {
	// create gorm postgres connection
	dsn := "host=localhost user=postgres dbname=autograd port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	conn, err := db.DB()
	if err != nil {
		panic(err)
	}

	conn.SetMaxIdleConns(8)
	conn.SetMaxOpenConns(100)

	return db
}

func DBTxFromGorm(tx *gorm.DB) (xsqlc.DBTX, bool) {
	dbtx, ok := tx.Statement.ConnPool.(*sql.Tx)
	return dbtx, ok
}
