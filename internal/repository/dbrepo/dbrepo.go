package dbrepo

import (
	"database/sql"
	"github.com/yaji1122/bookings-go/internal/repository"
)

type mysqlDBRepo struct {
	Pool *sql.DB
}

func NewMysqlRepo(pool *sql.DB) repository.DatabaseRepo {
	return &mysqlDBRepo{
		Pool: pool,
	}
}
