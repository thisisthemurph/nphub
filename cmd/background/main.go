package main

import (
	"github.com/joho/godotenv"
	"log/slog"
	"nphud/internal/config"
	"nphud/internal/repository"
	"nphud/internal/service"
	"nphud/pkg/np"
	"nphud/pkg/store"
	"os"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}

func main() {
	app := config.NewAppConfig(os.Getenv)
	database, err := store.GetOrCreate(app.Database.FullPath)
	if err != nil {
		panic(err)
	}
	defer database.Close()

	gameRepo := repository.NewGameRepository(database)
	snapshotRepo := repository.NewSnapshotRepository(database)
	snapshotFileService := service.NewSnapshotFileService(app.SnapshotBasePath, gameRepo, snapshotRepo)

	games, err := gameRepo.List()
	if err != nil {
		panic(err)
	}

	// Get a new snapshot for each of the games.
	for _, g := range games {
		slog.Debug("Background: fetching snapshot", "game", g)
		game := np.New(g.Number, g.APIKey)

		snapshotBytes, err := game.GetCurrentSnapshot()
		if err != nil {
			slog.Error("Error fetching new snapshot", "game", g)
			continue
		}

		_, err = snapshotFileService.Create(g.Number, g.APIKey, snapshotBytes)
		if err != nil {
			slog.Error("Error creating new snapshot", "game", g)
			continue
		}
	}
}
