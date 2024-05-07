package repository

import (
	"database/sql"
	"errors"
)

var ErrGameExists = errors.New("game already exists")

type GameRepository struct {
	db *sql.DB
}

func NewGameRepository(db *sql.DB) GameRepository {
	return GameRepository{
		db: db,
	}
}

// List all games in the games table.
func (gr GameRepository) List() ([]GameRow, error) {
	results := make([]GameRow, 0)

	stmt := "SELECT id, number, api_key FROM games;"
	rows, err := gr.db.Query(stmt)
	if err != nil {
		return results, err
	}
	defer rows.Close()

	for rows.Next() {
		var game GameRow
		if err := rows.Scan(&game.ID, &game.Number, &game.APIKey); err != nil {
			return results, err
		}

		results = append(results, game)
	}
	return results, nil
}

func (gr GameRepository) GetByNumberAndApiKey(number string, apiKey string) (*GameRow, error) {
	result := new(GameRow)
	stmt := "SELECT id, number, api_key FROM games WHERE number = ? AND api_key = ?;"
	if err := gr.db.QueryRow(stmt, number, apiKey).Scan(&result.ID, &result.Number, &result.APIKey); err != nil {
		return result, err
	}
	return result, nil
}

// Create a new row in the games table, returning the ID of the newly created record.
func (gr GameRepository) Create(number, apiKey string) (int64, error) {
	exists, err := gr.Exists(number, apiKey)
	if err != nil {
		return 0, err
	}

	if exists {
		return 0, ErrGameExists
	}

	stmt := "INSERT INTO games (number, api_key) VALUES (?, ?) returning id;"
	res, err := gr.db.Exec(stmt, number, apiKey)
	if err != nil {
		return 0, err
	}

	lastInsertID, err := res.LastInsertId()
	return lastInsertID, err
}

// Exists returns true if a record exists in the games table with the given game number and API key.
func (gr GameRepository) Exists(number, apiKey string) (bool, error) {
	stmt := "SELECT count(*) FROM games WHERE number = ? AND api_key = ?;"
	var count int
	err := gr.db.QueryRow(stmt, number, apiKey).Scan(&count)
	return count > 0, err
}
