package main

import (
	"database/sql"
	"nphud/internal/background/types"
	"nphud/internal/shared/model"
	"nphud/pkg/np"
	"testing"
	"time"
)

func getGamesStub() ([]types.GameWithSnapshots, error) {
	var games []types.GameWithSnapshots
	games = append(games, types.GameWithSnapshots{
		GameRow: model.GameRow{
			ID:             1,
			Number:         "123",
			APIKey:         "abc",
			PlayerUID:      1,
			StartTime:      time.Now().Add(-(time.Duration(24*7) * time.Hour)).UnixMilli(),
			TickRate:       60,
			ProductionRate: 24,
		},
		Snapshots: []model.SnapshotRow{{
			ID:        sql.NullInt64{},
			Path:      "/snapshots/file.json",
			CreatedAt: time.Now().Add(-(time.Duration(24*6) * time.Hour)).UnixMilli(),
		}},
	})
	return games, nil
}

func createSnapshotStub(_ string, _ []byte) (string, error) {
	return "", nil
}

func takeSnapshotStub(_ np.NeptunesPrideGame) ([]byte, error) {
	return make([]byte, 0), nil
}

func updateDatabaseStub(_, _, _ string) error {
	return nil
}

func TestRun(t *testing.T) {
	results, err := run(getGamesStub, takeSnapshotStub, createSnapshotStub, updateDatabaseStub)
	if err != nil {
		t.Error(err)
	}

	if results.failCount > 0 {
		t.Errorf("fail count: %d", results.failCount)
	}

	if results.snapshotsSkipped > 0 {
		t.Errorf("snapshots skipped: %d", results.snapshotsSkipped)
	}

	if results.snapshotsTaken != 1 {
		t.Errorf("snapshots taken: %d", results.snapshotsTaken)
	}
}

func TestGetLastTickTime(t *testing.T) {
	testCases := []struct {
		name                 string
		tickRate             int
		startTime            time.Time
		currentTime          time.Time
		expectedLastTickTime time.Time
	}{
		{
			name:                 "With60MinuteTicks",
			tickRate:             60,
			startTime:            time.Date(2023, time.December, 1, 0, 0, 0, 0, time.UTC),
			currentTime:          time.Date(2024, time.May, 12, 5, 47, 0, 0, time.UTC),
			expectedLastTickTime: time.Date(2024, time.May, 12, 5, 0, 0, 0, time.UTC),
		}, {
			name:                 "With30MinuteTicks",
			tickRate:             30,
			startTime:            time.Date(2023, time.December, 1, 0, 0, 0, 0, time.UTC),
			currentTime:          time.Date(2024, time.May, 12, 5, 12, 0, 0, time.UTC),
			expectedLastTickTime: time.Date(2024, time.May, 12, 5, 0, 0, 0, time.UTC),
		}, {
			name:                 "With15MinuteTicks",
			tickRate:             15,
			startTime:            time.Date(2023, time.December, 1, 0, 0, 0, 0, time.UTC),
			currentTime:          time.Date(2024, time.May, 12, 5, 12, 0, 0, time.UTC),
			expectedLastTickTime: time.Date(2024, time.May, 12, 5, 0, 0, 0, time.UTC),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lastTickTime, err := GetLastTickTime(tc.currentTime, tc.startTime, tc.tickRate)
			if err != nil {
				t.Error(err)
			}

			if lastTickTime != tc.expectedLastTickTime {
				t.Errorf("NextTickTime (%v) does not match expected (%v)", lastTickTime, tc.expectedLastTickTime)
			}
		})
	}
}

func TestGetNextTickTime(t *testing.T) {
	testCases := []struct {
		name                 string
		tickRate             int
		startTime            time.Time
		currentTime          time.Time
		expectedNextTickTime time.Time
	}{
		{
			name:                 "With60MinuteTicks",
			tickRate:             60,
			startTime:            time.Date(2023, time.December, 1, 0, 0, 0, 0, time.UTC),
			currentTime:          time.Date(2024, time.May, 12, 5, 47, 0, 0, time.UTC),
			expectedNextTickTime: time.Date(2024, time.May, 12, 6, 0, 0, 0, time.UTC),
		}, {
			name:                 "With30MinuteTicks",
			tickRate:             30,
			startTime:            time.Date(2023, time.December, 1, 0, 0, 0, 0, time.UTC),
			currentTime:          time.Date(2024, time.May, 12, 5, 12, 0, 0, time.UTC),
			expectedNextTickTime: time.Date(2024, time.May, 12, 5, 30, 0, 0, time.UTC),
		}, {
			name:                 "With15MinuteTicks",
			tickRate:             15,
			startTime:            time.Date(2023, time.December, 1, 0, 0, 0, 0, time.UTC),
			currentTime:          time.Date(2024, time.May, 12, 5, 12, 0, 0, time.UTC),
			expectedNextTickTime: time.Date(2024, time.May, 12, 5, 15, 0, 0, time.UTC),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			nextTickTime, err := GetNextTickTime(tc.currentTime, tc.startTime, tc.tickRate)
			if err != nil {
				t.Error(err)
			}

			if nextTickTime != tc.expectedNextTickTime {
				t.Errorf("NextTickTime (%v) does not match expected (%v)", nextTickTime, tc.expectedNextTickTime)
			}
		})
	}
}

func TestShouldTakeNewSnapshot(t *testing.T) {
	game := types.GameWithSnapshots{
		GameRow: model.GameRow{
			ID:             1,
			Number:         "123",
			APIKey:         "abc",
			PlayerUID:      1,
			StartTime:      time.Date(2023, time.December, 1, 0, 0, 0, 0, time.UTC).UnixMilli(),
			TickRate:       60,
			ProductionRate: 24,
		},
		Snapshots: []model.SnapshotRow{
			{
				ID:        sql.NullInt64{},
				Path:      "",
				CreatedAt: time.Now().Add(-(time.Duration(24*6) * time.Hour)).UnixMilli(),
			},
		},
	}

	should, err := shouldTakeNewSnapshot(game)
	if err != nil {
		t.Error(err)
	}
	if !should {
		t.Error("shouldTakeNewSnapshot should return true")
	}
}
