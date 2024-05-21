package model

import (
	"time"

	"github.com/google/uuid"
)

type Game struct {
	ExternalId     uuid.UUID
	Name           string
	Number         string
	PlayerUID      int
	APIKey         string
	StartTime      time.Time
	TickRate       int
	ProductionRate int
	Started        bool
	Paused         bool
	GameOver       bool
}
