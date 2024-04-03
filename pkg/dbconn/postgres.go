package dbconn

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/fahmifan/autograd/pkg/logs"
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

func DBTxFromGorm(tx *gorm.DB) (xsqlc.DBTX, error) {
	dbtx, ok := tx.Statement.ConnPool.(*sql.Tx)
	if !ok {
		db, ok := tx.Statement.ConnPool.(*sql.DB)
		if !ok {
			return nil, errors.New("cast to *sql.DB")
		}
		return db, nil
	}

	return dbtx, nil
}

func SqlcTransaction(ctx context.Context, db *sql.DB, fn func(xsqlc.DBTX) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return logs.ErrWrapCtx(ctx, err, "SqlcTransaction: begin")
	}

	err = fn(tx)
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			return logs.ErrWrapCtx(ctx, fmt.Errorf("%w: %w", err, err2), "SqlcTransaction: rollback")
		}

		logs.ErrCtx(ctx, err, "SqlcTransaction: callback")
		return nil
	}

	if err = tx.Commit(); err != nil {
		return logs.ErrWrapCtx(ctx, err, "SqlcTransaction: commit")
	}

	return nil
}
