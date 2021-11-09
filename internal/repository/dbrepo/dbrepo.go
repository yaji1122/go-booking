package dbrepo

import (
	"database/sql"
	"github.com/yaji1122/bookings-go/internal/config"
	"github.com/yaji1122/bookings-go/internal/repository"
)

type mysqlDBRepo struct {
	App  *config.AppConfig
	Pool *sql.DB
}

func NewMysqlRepo(pool *sql.DB, config *config.AppConfig) repository.DatabaseRepo {
	return &mysqlDBRepo{
		App:  config,
		Pool: pool,
	}
}
