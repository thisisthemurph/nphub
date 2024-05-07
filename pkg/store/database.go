package store

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func GetOrCreate(databaseFullPath string) (*sql.DB, error) {
	return sql.Open("sqlite3", databaseFullPath)
}
