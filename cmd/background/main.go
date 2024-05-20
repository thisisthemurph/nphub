package main

import (
	"log/slog"
	"nphud/internal/background/service"
	shared_service "nphud/internal/shared/service"
	"nphud/pkg/config"
	"nphud/pkg/np"
	"nphud/pkg/store"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	app := config.NewAppConfig(os.Getenv)
	database, err := store.GetOrCreate(app.Database.FullPath)
	if err != nil {
		panic(err)
	}
	defer database.Close()

	bgService := service.NewBackgroundService(database)
	ssFileService := shared_service.NewSnapshotFileService(app.SnapshotBasePath)

	if _, err = run(
		bgService.GetGamesRequiringNewSnapshot,
		takeSnapshot,
		ssFileService.Save,
		bgService.AddSnapshotToDatabase,
	); err != nil {
		slog.Error("could not fetch latest snapshots", "err", err)
	}
}

type runResult struct {
	snapshotsTaken int
	failCount      int
}

func run(
	getGames func() ([]np.NeptunesPrideGame, error),
	takeSnapshot func(game np.NeptunesPrideGame) ([]byte, error),
	createSnapshot func(string, []byte) (string, error),
	addSnapshotToDatabase func(string, string, string) error,
) (runResult, error) {
	var results runResult

	games, err := getGames()
	if err != nil {
		return results, err
	}

	for _, g := range games {
		snapshotBytes, err := takeSnapshot(g)
		if err != nil {
			slog.Error("could not take new snapshot", "game", g, "err", err)
			results.failCount += 1
			continue
		}

		snapshotFileName, err := createSnapshot(g.Number, snapshotBytes)
		if err != nil {
			slog.Error("could not create snapshot file", "game", g, "err", err)
			results.failCount += 1
			continue
		}

		if err = addSnapshotToDatabase(g.Number, g.APIKey, snapshotFileName); err != nil {
			slog.Error("could not add snapshot data to database", "game", g, "err", err)
			results.failCount += 1
		}

		results.snapshotsTaken += 1
	}

	slog.Info("snapshots collected", "results", results)
	return results, nil
}

// takeSnapshot executes the TakeSnapshot method on the game struct calling the
// NP API and retrieving snapshot data.
func takeSnapshot(game np.NeptunesPrideGame) ([]byte, error) {
	return game.TakeSnapshot()
}
