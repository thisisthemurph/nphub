package main

import (
	"errors"
	"github.com/joho/godotenv"
	"log/slog"
	bgservice "nphud/internal/background/service"
	"nphud/internal/background/types"
	"nphud/internal/shared/service"
	"nphud/pkg/config"
	"nphud/pkg/np"
	"nphud/pkg/store"
	"os"
	"time"
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

	bgService := bgservice.NewBackgroundService(database)
	ssFileService := service.NewSnapshotFileService(app.SnapshotBasePath)

	if _, err = run(
		bgService.ListGames,
		takeSnapshot,
		ssFileService.Save,
		bgService.AddSnapshotToDatabase,
	); err != nil {
		slog.Error("Error fetching latest snapshots", "err", err)
	}
}

type runResult struct {
	snapshotsTaken   int
	snapshotsSkipped int
	failCount        int
}

func run(
	getGames func() ([]types.GameWithSnapshots, error),
	takeSnapshot func(game np.NeptunesPrideGame) ([]byte, error),
	createSnapshot func(string, []byte) (string, error),
	updateDatabase func(string, string, string) error,
) (runResult, error) {
	var results runResult

	games, err := getGames()
	if err != nil {
		return results, err
	}

	for _, g := range games {
		shouldTakeSnapshot, err := shouldTakeNewSnapshot(g)
		if err != nil {
			slog.Error("Cannot determine if snapshot required", "err", err)
			results.failCount += 1
			continue
		}

		if !shouldTakeSnapshot {
			results.snapshotsSkipped += 1
			continue
		}

		slog.Debug("Background: fetching snapshot", "game", g)
		game := np.New(g.Number, g.APIKey)
		snapshotBytes, err := takeSnapshot(game)
		if err != nil {
			slog.Error("Error fetching new snapshot", "game", g)
			results.failCount += 1
			continue
		}

		snapshotFileName, err := createSnapshot(g.Number, snapshotBytes)
		if err != nil {
			slog.Error("Error creating new snapshot", "game", g)
			results.failCount += 1
			continue
		}

		if err = updateDatabase(g.Number, g.APIKey, snapshotFileName); err != nil {
			slog.Error("Error updating database", "game", g, "err", err)
			results.failCount += 1
		}

		results.snapshotsTaken += 1
	}

	return results, nil
}

// takeSnapshot executes the TakeSnapshot method on the game struct calling the
// NP API and retrieving snapshot data.
func takeSnapshot(game np.NeptunesPrideGame) ([]byte, error) {
	return game.TakeSnapshot()
}

// shouldTakeNewSnapshot determines if a snapshot is required for the given game.
func shouldTakeNewSnapshot(g types.GameWithSnapshots) (bool, error) {
	gameStartTime := time.Unix(0, g.StartTime*int64(time.Millisecond))
	lastTickTime, err := GetLastTickTime(time.Now(), gameStartTime, g.TickRate)
	if err != nil {
		return false, err
	}

	latestSnapshot, err := g.LatestSnapshot()
	if err != nil {
		slog.Debug("Error fetching latest snapshot", "err", err)
		return false, err
	}

	latestSnapshotCreatedAtTime := time.Unix(0, latestSnapshot.CreatedAt*int64(time.Millisecond))
	if latestSnapshotCreatedAtTime.After(lastTickTime) {
		slog.Debug("Skipping game", "startTime", gameStartTime, "lastTickTime", lastTickTime)
		return false, nil
	}

	return true, nil
}

// GetLastTickTime calculates the last tick time for a game.
// Parameters:
//
//	currentTime: the time from which to calculate, the result will be the time of the last tick prior to this value
//	startTime: the time the game started
//	tickRate: the tick rate of the game
//
// Returns an error if there is a possible divide by zero issue.
func GetLastTickTime(currentTime, startTime time.Time, tickRate int) (time.Time, error) {
	// Calculate the number of ticks between now and the start
	timeDiff := currentTime.Sub(startTime)
	minutesDiff := int(timeDiff.Minutes())

	if minutesDiff == 0 || tickRate == 0 {
		return time.Time{}, errors.New("cannot divide by zero")
	}
	elapsedTicks := minutesDiff / tickRate

	return startTime.Add(time.Duration(elapsedTicks*tickRate) * time.Minute), nil
}

// GetNextTickTime calculates the next tick time for a game.
// Parameters:
//
//	currentTime: the time from which to calculate, the result will be the time of the tick following to this value
//	startTime: the time the game started
//	tickRate: the tick rate of the game
//
// Returns an error if there is a possible divide by zero issue.
func GetNextTickTime(currentTime, startTime time.Time, tickRate int) (time.Time, error) {
	lastTickTime, err := GetLastTickTime(currentTime, startTime, tickRate)
	if err != nil {
		return time.Time{}, err
	}

	return lastTickTime.Add(time.Duration(tickRate) * time.Minute), nil
}
