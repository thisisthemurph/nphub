package model

import "time"

type Game struct {
	Name           string
	Number         string
	PlayerUID      string
	APIKey         string
	StartTime      time.Time
	TickRate       int
	ProductionRate int
	Started        bool
	Paused         bool
	GameOver       bool
}
