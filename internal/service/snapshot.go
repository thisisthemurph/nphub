package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"nphud/internal/repository"
	"os"
	"path/filepath"
)

var ErrSnapshotDataMalformed = errors.New("malformed snapshot data")

// SnapshotFileService is a service for dealing with snapshot files.
type SnapshotFileService struct {
	basePath     string
	gameRepo     repository.GameRepository
	snapshotRepo repository.SnapshotRepository
}

// NewSnapshotFileService creates a new SnapshotFileService.
func NewSnapshotFileService(
	basePath string,
	gameRepo repository.GameRepository,
	snapshotRepo repository.SnapshotRepository,
) SnapshotFileService {
	return SnapshotFileService{
		basePath:     basePath,
		gameRepo:     gameRepo,
		snapshotRepo: snapshotRepo,
	}
}

// Create a new snapshot file and return the name of the newly created file.
func (s SnapshotFileService) Create(gameNumber, apiKey string, data []byte) (string, error) {
	fileName, err := s.makeFileName(gameNumber, data)
	if err != nil {
		return fileName, err
	}

	gameRow, err := s.gameRepo.GetByNumberAndApiKey(gameNumber, apiKey)
	if err != nil {
		return fileName, err
	}

	if err = os.WriteFile(filepath.Join(s.basePath, fileName), data, 0666); err != nil {
		return fileName, err
	}

	if err = s.snapshotRepo.Create(gameRow.ID, fileName); err != nil {
		slog.Error("Error inserting snapshot data into database", "game", gameRow)
		return fileName, err
	}

	return fileName, err
}

// makeFileName creates a file name with details present in the provided bytes JSON data.
func (s SnapshotFileService) makeFileName(gameNumber string, data []byte) (string, error) {
	type FileNameMetadata struct {
		Metadata struct {
			Now       int64 `json:"now"`
			PlayerUID int   `json:"player_uid"`
			Tick      int   `json:"tick"`
		} `json:"scanning_data"`
	}

	var gameMetadata FileNameMetadata
	if err := json.Unmarshal(data, &gameMetadata); err != nil {
		slog.Error("Error unmarshaling game metadata", "err", err)
		return "", ErrSnapshotDataMalformed
	}
	md := gameMetadata.Metadata
	return fmt.Sprintf("%s_%v_%v_%v.json", gameNumber, md.Tick, md.Now, md.PlayerUID), nil
}
