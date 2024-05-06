package repository

import (
	"database/sql"
	"errors"
	"fmt"
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

		fmt.Printf("%+v\n", game)
		results = append(results, game)
	}
	return results, nil
}

func (gr GameRepository) Create(number, apiKey string) error {
	exists, err := gr.Exists(number, apiKey)
	if err != nil {
		return err
	}

	if exists {
		return ErrGameExists
	}

	stmt := "INSERT INTO games (number, api_key) VALUES (?, ?)"
	_, err = gr.db.Exec(stmt, number, apiKey)
	return err
}

func (gr GameRepository) Exists(number, apiKey string) (bool, error) {
	stmt := "SELECT count(*) FROM games WHERE number = ? AND api_key = ?;"
	var count int
	err := gr.db.QueryRow(stmt, number, apiKey).Scan(&count)
	return count > 0, err
}
