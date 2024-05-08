package handler

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"nphud/internal/repository"
	"nphud/internal/service"
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
	gameRepo            repository.GameRepository
	snapshotRepo        repository.SnapshotRepository
	snapshotFileService service.SnapshotFileService
}

func NewGameHandler(
	gameRepo repository.GameRepository,
	snapshotRepo repository.SnapshotRepository,
	snapshotFileService service.SnapshotFileService,
) GameHandler {
	return GameHandler{
		gameRepo:            gameRepo,
		snapshotRepo:        snapshotRepo,
		snapshotFileService: snapshotFileService,
	}
}

type GameResponse struct {
	Number string `json:"number"`
	APIKey string `json:"api_key"`
}

type GameListResponse struct {
	Count int            `json:"count"`
	Games []GameResponse `json:"games"`
}

func (h GameHandler) ListGames(c echo.Context) error {
	games, err := h.gameRepo.List()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, JSONError{Message: "Internal Server Error"})
	}

	response := GameListResponse{Count: len(games)}
	for _, g := range games {
		response.Games = append(response.Games, GameResponse{Number: g.Number, APIKey: g.APIKey})
	}

	return c.JSON(http.StatusOK, response)
}

func (h GameHandler) CreateNewGame(c echo.Context) error {
	req := new(NewGameRequest)
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, JSONError{Message: "Invalid JSON"})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, JSONError{Message: "Invalid request"})
	}

	// Create a new game instance and get the current snapshot.
	game := np.New(req.GameNumber, req.APIKey)
	snapshotBytes, err := game.GetCurrentSnapshot()
	if err != nil {
		return c.JSON(http.StatusBadRequest, JSONError{Message: err.Error()})
	}

	// Insert the game into the games row.
	_, err = h.gameRepo.Create(req.GameNumber, req.APIKey)
	if err != nil {
		if errors.Is(err, repository.ErrGameExists) {
			return c.JSON(http.StatusConflict, JSONError{Message: err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, JSONError{Message: "Internal Server Error"})
	}

	// Save the snapshot file and create a record in the snapshots table.
	_, err = h.snapshotFileService.Create(req.GameNumber, req.APIKey, snapshotBytes)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, JSONError{Message: "Internal Server Error"})
	}

	return c.JSON(http.StatusCreated, nil)
}

func (h GameHandler) GetGame(c echo.Context) error {
	req := new(NewGameRequest)
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, JSONError{Message: "Invalid JSON"})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, JSONError{Message: "Invalid request"})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, JSONError{Message: "Invalid request"})
	}

	err := h.snapshotFileService.GetMostRecent(req.GameNumber, req.APIKey)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, JSONError{Message: "Internal Server Error"})
	}
	return nil
}
