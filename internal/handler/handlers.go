package handler

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"nphud/internal/repository"
	"nphud/internal/services"
	"nphud/pkg/np"
)

type NewGameRequest struct {
	GameNumber string `json:"game_number" validate:"required"`
	APIKey     string `json:"api_key" validate:"required"`
}

type JSONError struct {
	Message string `json:"error"`
}

type GameHandler struct {
	repo                repository.GameRepository
	snapshotFileService services.SnapshotFileService
}

func NewGameHandler(repo repository.GameRepository, snapshotFileService services.SnapshotFileService) GameHandler {
	return GameHandler{
		repo:                repo,
		snapshotFileService: snapshotFileService,
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
	snapshotBytes, err := game.GetCurrentSnapshot()
	if err != nil {
		return c.JSON(http.StatusBadRequest, JSONError{Message: err.Error()})
	}

	if err = h.repo.Create(req.GameNumber, req.APIKey); err != nil {
		if errors.Is(err, repository.ErrGameExists) {
			return c.JSON(http.StatusConflict, JSONError{Message: err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, JSONError{Message: "Internal Server Error"})
	}

	if err = h.snapshotFileService.Create(game.Number, snapshotBytes); err != nil {
		return c.JSON(http.StatusInternalServerError, JSONError{Message: "Internal Server Error"})
	}

	return c.JSON(http.StatusCreated, nil)
}
