package service

import (
	"database/sql"
	"errors"
	"log/slog"
	"nphud/pkg/np"
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

func (bs BackgroundService) GetGamesRequiringNewSnapshot() ([]np.NeptunesPrideGame, error) {
	stmt := `
		select number, api_key
		from games
		where next_snapshot_at <= ?;`

	rows, err := bs.db.Query(stmt, time.Now().UnixMilli())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	games := make([]np.NeptunesPrideGame, 0)
	for rows.Next() {
		var number string
		var key string
		if err := rows.Scan(&number, &key); err != nil {
			return nil, err
		}
		games = append(games, np.New(number, key))
	}

	return games, nil
}

func (bs BackgroundService) AddSnapshotToDatabase(gameNumber, apiKey, snapshotFileName string) error {
	tx, err := bs.db.Begin()
	if err != nil {
		slog.Error("could not begin transaction", "err", err)
		return err
	}
	defer tx.Rollback()

	var (
		gameRowId       int64
		startTimeMillis int64
		tickRate        int
	)

	stmt := "select id, start_time, tick_rate from games where number = ? and api_key = ?;"
	err = tx.QueryRow(stmt, gameNumber, apiKey).Scan(&gameRowId, &startTimeMillis, &tickRate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrGameNotFound
		}
		return err
	}

	gameStartTime := time.Unix(0, startTimeMillis*int64(time.Millisecond))
	nextTickTime, err := np.CalculateNextTickTime(gameStartTime, tickRate)
	if err != nil {
		slog.Error("error calculating next tick time", "err", err)
		return err
	}

	stmt = "update games set next_snapshot_at = ? where id = ?;"
	if _, err := tx.Exec(stmt, nextTickTime.UnixMilli(), gameRowId); err != nil {
		slog.Error("error updating next_snapshot_time", "err", err)
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
