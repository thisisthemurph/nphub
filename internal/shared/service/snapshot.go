package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"nphud/pkg/np"
	"os"
	"path/filepath"
)

var (
	ErrFileNotFound          = errors.New("snapshot file does not exist")
	ErrSnapshotDataMalformed = errors.New("malformed snapshot data")
)

// SnapshotFileService is a service for dealing with snapshot files.
type SnapshotFileService struct {
	basePath string
}

// NewSnapshotFileService creates a new SnapshotFileService.
func NewSnapshotFileService(basePath string) SnapshotFileService {
	return SnapshotFileService{
		basePath: basePath,
	}
}

// Save persists the given data bytes as a snapshot JSON file.
func (s SnapshotFileService) Save(gameNumber string, data []byte) (string, error) {
	fileName, err := s.makeFileName(gameNumber, data)
	if err != nil {
		return fileName, err
	}

	if err = os.WriteFile(filepath.Join(s.basePath, fileName), data, 0666); err != nil {
		return fileName, err
	}

	return fileName, nil
}

// Get returns snapshot data for the given fileName.
// Returns an error if the file does not exist or if the JSON cannot be unmarshalled.
func (s SnapshotFileService) Get(fileName string) (np.APIResponse, error) {
	var snapshot np.APIResponse
	filePath := fmt.Sprintf("%s/%s", s.basePath, fileName)

	b, err := os.ReadFile(filePath)
	if err != nil {
		return snapshot, err
	}

	if err = json.Unmarshal(b, &snapshot); err != nil {
		return snapshot, err
	}

	return snapshot, nil
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
		slog.Error("Error unmarshalling game metadata", "err", err)
		return "", ErrSnapshotDataMalformed
	}
	md := gameMetadata.Metadata
	return fmt.Sprintf("%s_%v_%v_%v.json", gameNumber, md.Tick, md.Now, md.PlayerUID), nil
}
