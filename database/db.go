package database

import (
	"database/sql"
)

func InitDB(driver, dbConnectionString string) (*sql.DB, error) {
	return sql.Open(driver, dbConnectionString)
}
