package command

import (
	"context"
	"database/sql"
	"nphud/internal/shared/service"
	"nphud/pkg/np/model"
)

type GetLatestScanningDataQuery struct {
	GameNumber string
	PlayerUID  int
}

func NewGetLatestScanningDataQuery(gameNumber string, playerUID int) *GetLatestScanningDataQuery {
	return &GetLatestScanningDataQuery{
		GameNumber: gameNumber,
		PlayerUID:  playerUID,
	}
}

type GetLatestScanningDataQueryHandler struct {
	db                  *sql.DB
	snapshotFileService service.SnapshotFileService
}

func NewGetLatestScanningDataQueryHandler(db *sql.DB, snapshotFileService service.SnapshotFileService) *GetLatestScanningDataQueryHandler {
	return &GetLatestScanningDataQueryHandler{
		db:                  db,
		snapshotFileService: snapshotFileService,
	}
}

func (h *GetLatestScanningDataQueryHandler) Handle(ctx context.Context, cmd *GetLatestScanningDataQuery) (model.ScanningData, error) {
	stmt := `
		select s.path
		from games g
		join snapshots s on g.id = s.game_id
		where g.number = ? and g.player_uid = ?
		order by created_at desc
		limit 1;`

	var snapshotFileName string
	err := h.db.QueryRow(stmt, cmd.GameNumber, cmd.PlayerUID).Scan(&snapshotFileName)
	if err != nil {
		return model.ScanningData{}, err
	}

	snapshot, err := h.snapshotFileService.Get(snapshotFileName)
	if err != nil {
		return model.ScanningData{}, err
	}

	return snapshot.ScanningData, nil
}
