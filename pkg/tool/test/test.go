package test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func GetSQLMock() (*sqlx.DB, sqlmock.Sqlmock) {
	db, mock, _ := sqlmock.New(sqlmock.MonitorPingsOption(true))
	return sqlx.NewDb(db, "mysql"), mock
}
