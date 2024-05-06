package store

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func GetOrCreate() (*sql.DB, error) {
	database, _ := sql.Open("sqlite3", "./np.db")
	stmt, err := database.Prepare(`
		CREATE TABLE IF NOT EXISTS games (
		    id INTEGER PRIMARY KEY,
		    number TEXT,
		    api_key TEXT
		);
	`)
	if err != nil {
		return nil, err
	}
	if _, err = stmt.Exec(); err != nil {
		return nil, err
	}
	return database, nil
}
