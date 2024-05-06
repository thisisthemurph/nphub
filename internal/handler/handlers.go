package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"nphud/internal/repository"
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
	repo repository.GameRepository
}

func NewGameHandler(repo repository.GameRepository) GameHandler {
	return GameHandler{
		repo: repo,
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

	if err = h.repo.Create(req.GameNumber, req.APIKey); err != nil {
		if errors.Is(err, repository.ErrGameExists) {
			return c.JSON(http.StatusConflict, JSONError{Message: err.Error()})
		}
		return echo.NewHTTPError(http.StatusInternalServerError, JSONError{Message: "Internal Server Error"})
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
