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
	gameRowID, err := h.gameRepo.Create(req.GameNumber, req.APIKey)
	if err != nil {
		if errors.Is(err, repository.ErrGameExists) {
			return c.JSON(http.StatusConflict, JSONError{Message: err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, JSONError{Message: "Internal Server Error"})
	}

	// Save the snapshot file.
	snapshotFileName, err := h.snapshotFileService.Create(game.Number, snapshotBytes)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, JSONError{Message: "Internal Server Error"})
	}

	// Insert an instance into the snapshots.
	if err = h.snapshotRepo.Create(gameRowID, snapshotFileName); err != nil {
		return c.JSON(http.StatusInternalServerError, JSONError{Message: "Internal Server Error"})
	}

	return c.JSON(http.StatusCreated, nil)
}
