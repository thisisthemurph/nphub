package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"nphud/pkg/np"
	"os"
)

type GameMetadata struct {
	Metadata struct {
		Now       int64 `json:"now"`
		PlayerUID int   `json:"player_uid"`
		Tick      int   `json:"tick"`
	} `json:"scanning_data"`
}

type NewGameRequest struct {
	GameNumber string `json:"game_number" validate:"required"`
	APIKey     string `json:"api_key" validate:"required"`
}

type JSONError struct {
	Message string `json:"error"`
}

type GameHandler struct {
	db *sql.DB
}

func NewGameHandler(db *sql.DB) GameHandler {
	return GameHandler{
		db: db,
	}
}

func (h GameHandler) CreateNewGame(c echo.Context) error {
	req := new(NewGameRequest)
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, JSONError{Message: "Invalid JSON"})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, JSONError{Message: "Invalid request"})
	}

	game := np.New(req.GameNumber, req.APIKey)
	snapshot, err := game.GetCurrentSnapshot()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, JSONError{Message: err.Error()})
	}

	stmt := "INSERT INTO games (number, api_key) VALUES (?, ?)"
	_, err = h.db.Exec(stmt, req.GameNumber, req.APIKey)
	if err != nil {
		return err
	}

	var gameMetadata GameMetadata
	if err := json.Unmarshal(snapshot, &gameMetadata); err != nil {
		return err
	}
	md := gameMetadata.Metadata
	fileName := fmt.Sprintf("snapshots/%s_%v_%v_%v.json", req.GameNumber, md.Tick, md.Now, md.PlayerUID)
	if err := os.WriteFile(fileName, snapshot, 0666); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, nil)
}
