package model

import "time"

type Game struct {
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
