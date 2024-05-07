package main

import (
	"fmt"
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

	snapshotFileService := service.NewSnapshotFileService(app.SnapshotBasePath)
	gameRepo := repository.NewGameRepository(database)
	snapshotRepo := repository.NewSnapshotRepository(database)

	games, err := gameRepo.List()
	if err != nil {
		panic(err)
	}

	// Get a new snapshot for each of the games.
	for _, g := range games {
		fmt.Printf("Game ID: %d\n", g.ID)
		game := np.New(g.Number, g.APIKey)
		snapshotBytes, err := game.GetCurrentSnapshot()
		if err != nil {
			slog.Error("Error fetching new snapshot", "game", g)
			continue
		}

		fileName, err := snapshotFileService.Create(g.Number, snapshotBytes)
		if err != nil {
			slog.Error("Error creating new snapshot", "game", g)
			continue
		}

		if err = snapshotRepo.Create(g.ID, fileName); err != nil {
			slog.Error("Error inserting snapshot data into database", "game", g)
			continue
		}
	}
}
