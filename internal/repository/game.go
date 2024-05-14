package repository

import (
	"database/sql"
	"errors"
)

var ErrGameExists = errors.New("game already exists")

type GameRepository struct {
	db *sql.DB
}

func NewGameRepository(db *sql.DB) GameRepository {
	return GameRepository{
		db: db,
	}
}

type GameWithSnapshots struct {
	GameRow
	Snapshots []SnapshotRow
}

func (gws GameWithSnapshots) LatestSnapshot() (SnapshotRow, error) {
	if gws.Snapshots == nil || len(gws.Snapshots) == 0 {
		return SnapshotRow{}, errors.New("no game snapshots found")
	}

	var latestSnapshot SnapshotRow
	for _, snapshot := range gws.Snapshots {
		if snapshot.CreatedAt > latestSnapshot.CreatedAt {
			latestSnapshot = snapshot
		}
	}
	return latestSnapshot, nil
}

// List all games in the games table.
func (gr GameRepository) List() ([]GameWithSnapshots, error) {
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

	rows, err := gr.db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	gamesMap := make(map[int64]GameWithSnapshots)
	for rows.Next() {
		var g GameRow
		var s SnapshotRow

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
				gamesMap[g.ID] = GameWithSnapshots{
					GameRow:   g,
					Snapshots: []SnapshotRow{s},
				}
			} else {
				gamesMap[g.ID] = GameWithSnapshots{
					GameRow: g,
				}
			}
		}
	}

	var gamesWithSnapshots []GameWithSnapshots
	for _, g := range gamesMap {
		gamesWithSnapshots = append(gamesWithSnapshots, g)
	}
	return gamesWithSnapshots, nil
}

func (gr GameRepository) GetByNumberAndApiKey(number string, apiKey string) (*GameRow, error) {
	result := new(GameRow)
	stmt := "SELECT id, number, api_key FROM games WHERE number = ? AND api_key = ?;"
	if err := gr.db.QueryRow(stmt, number, apiKey).Scan(&result.ID, &result.Number, &result.APIKey); err != nil {
		return result, err
	}
	return result, nil
}

// Create a new row in the games table, returning the ID of the newly created record.
func (gr GameRepository) Create(g GameRowCreate) (int64, error) {
	exists, err := gr.Exists(g.Number, g.APIKey)
	if err != nil {
		return 0, err
	}

	if exists {
		return 0, ErrGameExists
	}

	stmt := `
		insert into games (
			number,
			api_key,
			player_uid,
			start_time,
			tick_rate,
			production_rate
		) VALUES (?, ?, ?, ?, ?, ?);`

	res, err := gr.db.Exec(stmt, g.Number, g.APIKey, g.PlayerUID, g.StartTimeRaw, g.TickRate, g.ProductionRate)
	if err != nil {
		return 0, err
	}

	lastInsertID, err := res.LastInsertId()
	return lastInsertID, err
}

// Exists returns true if a record exists in the games table with the given game number and API key.
func (gr GameRepository) Exists(number, apiKey string) (bool, error) {
	stmt := "SELECT count(*) FROM games WHERE number = ? AND api_key = ?;"
	var count int
	err := gr.db.QueryRow(stmt, number, apiKey).Scan(&count)
	return count > 0, err
}
