package command

import (
	"context"
	"database/sql"
)

type GameExistsQuery struct {
	Number string
	APIKey string
}

func NewGameExistsQuery(number, apiKey string) *GameExistsQuery {
	return &GameExistsQuery{
		Number: number,
		APIKey: apiKey,
	}
}

type GameExistsQueryHandler struct {
	db *sql.DB
}

func NewGameExistsQueryHandler(db *sql.DB) *GameExistsQueryHandler {
	return &GameExistsQueryHandler{
		db: db,
	}
}

func (c *GameExistsQueryHandler) Handle(ctx context.Context, cmd *GameExistsQuery) (bool, error) {
	// Check if a game exists with the exact params

	stmt := "select count(*) from games where number = ? and api_key = ?;"

	var count int
	err := c.db.QueryRow(stmt, cmd.Number, cmd.APIKey).Scan(&count)
	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, nil
	}

	// Check if a game exists but this is a new API key

	return true, nil

}
