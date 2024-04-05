package dbconn

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/fahmifan/autograd/pkg/logs"
	"github.com/fahmifan/autograd/pkg/xsqlc"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zerologadapter"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func MustPostgres(debug bool) *gorm.DB {
	// create gorm postgres connection
	dsn := "host=localhost user=postgres dbname=autograd port=5432 sslmode=disable"
	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		panic(err)
	}

	conn.SetMaxIdleConns(2)
	conn.SetMaxOpenConns(10)

	if debug {
		loggerAdapter := zerologadapter.New(zerolog.New(os.Stdout))
		conn = sqldblogger.OpenDriver(dsn, conn.Driver(), loggerAdapter)
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: conn,
	}))
	if err != nil {
		panic(err)
	}

	return db
}

func NewSqlcFromGorm(tx *gorm.DB) (*xsqlc.Queries, error) {
	dbtx, err := DBTxFromGorm(tx)
	if err != nil {
		return nil, err
	}

	return xsqlc.New(dbtx), nil
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
