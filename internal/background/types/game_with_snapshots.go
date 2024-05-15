package types

import (
	"errors"
	"nphud/internal/shared/model"
)

var ErrGameHasNoSnapshots = errors.New("game has no snapshots")

type GameWithSnapshots struct {
	model.GameRow
	Snapshots []model.SnapshotRow
}

func (gws GameWithSnapshots) LatestSnapshot() (model.SnapshotRow, error) {
	if gws.Snapshots == nil || len(gws.Snapshots) == 0 {
		return model.SnapshotRow{}, ErrGameHasNoSnapshots
	}

	var latestSnapshot model.SnapshotRow
	for _, snapshot := range gws.Snapshots {
		if snapshot.CreatedAt > latestSnapshot.CreatedAt {
			latestSnapshot = snapshot
		}
	}
	return latestSnapshot, nil
}
