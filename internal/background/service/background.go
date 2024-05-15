package service

import (
	"database/sql"
	"errors"
	"log/slog"
	"nphud/internal/background/types"
	"nphud/internal/shared/model"
	"time"
)

var ErrGameNotFound = errors.New("game not found")

type BackgroundService struct {
	db *sql.DB
}

func NewBackgroundService(db *sql.DB) BackgroundService {
	return BackgroundService{
		db: db,
	}
}

func (bs BackgroundService) ListGames() ([]types.GameWithSnapshots, error) {
	stmt := `
	select
	  g.id,
	  g.number,
	  g.player_uid,
	  g.api_key,
	  g.start_time,
	  g.tick_rate,
	  g.production_rate,
	  s.id,
	  s.path,
	  s.created_at
	from games g
	left join snapshots s on s.game_id = g.id
	order by g.id, s.created_at;`

	rows, err := bs.db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	gamesMap := make(map[int64]types.GameWithSnapshots)
	for rows.Next() {
		var g model.GameRow
		var s model.SnapshotRow

		err = rows.Scan(
			&g.ID,
			&g.Number,
			&g.PlayerUID,
			&g.APIKey,
			&g.StartTime,
			&g.TickRate,
			&g.ProductionRate,
			&s.ID,
			&s.Path,
			&s.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if gws, ok := gamesMap[g.ID]; ok {
			gws.Snapshots = append(gws.Snapshots, s)
			gamesMap[g.ID] = gws
		} else {
			if s.ID.Valid {
				gamesMap[g.ID] = types.GameWithSnapshots{
					GameRow:   g,
					Snapshots: []model.SnapshotRow{s},
				}
			} else {
				gamesMap[g.ID] = types.GameWithSnapshots{
					GameRow: g,
				}
			}
		}
	}

	var gamesWithSnapshots []types.GameWithSnapshots
	for _, g := range gamesMap {
		gamesWithSnapshots = append(gamesWithSnapshots, g)
	}
	return gamesWithSnapshots, nil
}

func (bs BackgroundService) AddSnapshotToDatabase(gameNumber, apiKey, snapshotFileName string) error {
	tx, err := bs.db.Begin()
	if err != nil {
		slog.Error("could not begin transaction", "err", err)
		return err
	}
	defer tx.Rollback()

	stmt := "select id from games where number = ? and api_key = ?;"
	var gameRowId int64
	err = tx.QueryRow(stmt, gameNumber, apiKey).Scan(&gameRowId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrGameNotFound
		}
		return err
	}

	stmt = "insert into snapshots (game_id, path, created_at) values (?, ?, ?);"
	if _, err = tx.Exec(stmt, gameRowId, snapshotFileName, time.Now().UnixMilli()); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		slog.Error("could not commit transaction", "err", err)
		return err
	}

	return nil
}
