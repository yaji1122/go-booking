package driver

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type Pool struct {
	SQL *sql.DB
}

var dbConn = &Pool{}

const maxOpenConns = 10
const connMaxIdleTime = 5
const connMaxLifetime = 5 * time.Minute

func ConnectSQL(dsn string) (*Pool, error) {
	d, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}

	d.SetMaxOpenConns(maxOpenConns)
	d.SetMaxIdleConns(connMaxIdleTime)
	d.SetConnMaxLifetime(connMaxLifetime)

	dbConn.SQL = d
	if testDB(dbConn.SQL) != nil {
		panic(err)
	}
	return dbConn, err
}

func NewDatabase(dsn string) (*sql.DB, error) {
	pool, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func testDB(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		return err
	}
	return nil
}
