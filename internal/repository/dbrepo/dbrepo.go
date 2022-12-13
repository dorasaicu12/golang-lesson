package dbrepo

import (
	"database/sql"

	"github.com/dorasaicu12/booking/internal/config"
	"github.com/dorasaicu12/booking/internal/repository"
)

type mysqlDBrepo struct {
	App *config.AppConfig
	DB *sql.DB
}
func NewMysqlRepo (conn *sql.DB,a *config.AppConfig) repository.DatabaseRepo{
	return &mysqlDBrepo{
		App: a,
		DB:conn,
	}
}