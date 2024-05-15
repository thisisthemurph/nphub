package model

import "time"

type Game struct {
	Number         string
	PlayerUID      string
	APIKey         string
	StartTime      time.Time
	TickRate       int
	ProductionRate int
}
