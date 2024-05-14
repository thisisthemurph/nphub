package main

import (
	"github.com/labstack/echo/v4"
	"nphud/internal/handler"
	"nphud/internal/repository"
	"nphud/internal/service"
)

func makeRoutes(
	e *echo.Echo,
	gameRepository repository.GameRepository,
	snapshotRepository repository.SnapshotRepository,
	snapshotFileService service.SnapshotFileService,
) {
	gameHandler := handler.NewGameHandler(gameRepository, snapshotRepository, snapshotFileService)
	e.GET("/game", gameHandler.ListGames)
	e.POST("/game", gameHandler.CreateNewGame)

	e.GET("/snapshot", gameHandler.GetGame)
}
