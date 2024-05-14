package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"nphud/internal/config"
	"nphud/internal/repository"
	"nphud/internal/service"
	"nphud/pkg/store"
	"os"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}

func main() {
	app := config.NewAppConfig(os.Getenv)
	database, err := store.GetOrCreate(app.Database.FullPath)
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
	snapshotRepository := repository.NewSnapshotRepository(database)
	snapshotFileService := service.NewSnapshotFileService(app.SnapshotBasePath, gameRepository, snapshotRepository, database)

	makeRoutes(e, gameRepository, snapshotRepository, snapshotFileService)
	e.Logger.Fatal(e.Start(app.ListenAddress))
}
