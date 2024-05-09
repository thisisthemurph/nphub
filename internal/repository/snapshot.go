package repository

import (
	"database/sql"
	"time"
)

type SnapshotRepository struct {
	db *sql.DB
}

func NewSnapshotRepository(db *sql.DB) SnapshotRepository {
	return SnapshotRepository{
		db: db,
	}
}

// Create a new record in the snapshots table.
func (r SnapshotRepository) Create(gameRowID int64, fileName string) error {
	stmt := "insert into snapshots (game_id, path, created_at) values (?, ?, ?);"

	createdAt := time.Now().Unix()
	_, err := r.db.Exec(stmt, gameRowID, fileName, createdAt)
	return err
}
