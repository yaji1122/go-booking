package driver

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

const maxOpenConns = 10
const connMaxIdleTime = 5
const connMaxLifetime = 5 * time.Minute

//CreateDatabaseConnectionPool 建立連線池
func CreateDatabaseConnectionPool(dataSource string) *sql.DB {
	pool, err := sql.Open("mysql", dataSource)
	checkErr(err)

	pool.SetMaxOpenConns(maxOpenConns)
	pool.SetMaxIdleConns(connMaxIdleTime)
	pool.SetConnMaxLifetime(connMaxLifetime)
	//測試連線
	err = pool.Ping()
	checkErr(err)

	return pool
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
