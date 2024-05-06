package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"nphud/internal/handler"
	"nphud/internal/repository"
	"nphud/pkg/store"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	database, err := store.GetOrCreate()
	if err != nil {
		log.Fatal(err)
	}
	if err = database.Ping(); err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Validator = &CustomValidator{validator: validator.New()}

	gameRepository := repository.NewGameRepository(database)
	gameHandler := handler.NewGameHandler(gameRepository)
	e.POST("/game", gameHandler.CreateNewGame)

	e.Logger.Fatal(e.Start(":42069"))
}
