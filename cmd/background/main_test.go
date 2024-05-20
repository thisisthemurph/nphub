package main

import (
	"nphud/pkg/np"
	"testing"
)

func getGamesStub() ([]np.NeptunesPrideGame, error) {
	var games []np.NeptunesPrideGame
	games = append(games, np.New("123", "abc"))
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

	if results.snapshotsTaken != 1 {
		t.Errorf("snapshots taken: %d", results.snapshotsTaken)
	}
}
