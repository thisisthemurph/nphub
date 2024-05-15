package model

import "database/sql"

type GameRow struct {
	ID             int64
	Number         string
	APIKey         string
	PlayerUID      int
	StartTime      int64
	TickRate       int
	ProductionRate int
}

type SnapshotRow struct {
	ID        sql.NullInt64
	Path      string
	CreatedAt int64
}
