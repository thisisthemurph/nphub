package model

import (
	"encoding/json"
	"time"
)

type ScanningData struct {
	Fleets            map[int]Fleet `json:"fleets"`
	FleetSpeed        float64       `json:"fleet_speed"`
	Paused            bool          `json:"paused"`
	Productions       int           `json:"productions"`
	TickFragment      float64       `json:"tick_fragment"` // Percentage of current tick
	Now               time.Time     `json:"now"`
	TickRate          int           `json:"tick_rate"`          // Number of minutes per tick
	ProductionRate    int           `json:"production_rate"`    // Number of ticks per production cycle
	ProductionCounter int           `json:"production_counter"` // Current tick within the production cycle
	Stars             map[int]Star  `json:"stars"`
	StarsForVictory   int           `json:"stars_for_victory"`
	GameOver          bool          `json:"game_over"` // int in original JSON
	Started           bool          `json:"started"`
	StartTime         time.Time     `json:"start_time"`
	TotalStars        int           `json:"total_stars"`
	TradeScanned      int           `json:"trade_scanned"`
	Tick              int           `json:"tick"`
	TradeCost         int           `json:"trade_cost"`
	Name              string        `json:"name"`
	PlayerUID         int           `json:"player_uid"`
	Admin             bool          `json:"admin"`      // int in original JSON
	TurnBased         bool          `json:"turn_based"` // int in original JSON
	War               int           `json:"war"`        // Unknown purpose
	Players           PlayerList    `json:"players"`
	TurnBasedTimeOut  int           `json:"turn_based_time_out"`

	StartTimeRaw int64
}

// GetNextProductionTime returns the time for the next production cycle.
func (sd *ScanningData) GetNextProductionTime() time.Time {
	productionCycleMinutes := sd.ProductionRate * sd.TickRate
	minutesIntoCurrentProductionCycle := sd.ProductionCounter * sd.TickRate
	minutesUntilNextProductionCycle := productionCycleMinutes - minutesIntoCurrentProductionCycle
	return time.Now().Add(time.Duration(minutesUntilNextProductionCycle) * time.Minute)
}

func (sd *ScanningData) UnmarshalJSON(data []byte) error {
	type Alias ScanningData
	aux := struct {
		Now       int64 `json:"now"`
		StartTime int64 `json:"start_time"`
		GameOver  int   `json:"game_over"`
		Admin     int   `json:"admin"`
		TurnBased int   `json:"turn_based"`
		*Alias
	}{
		Alias: (*Alias)(sd),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	sd.Now = timeFromUnixMilliseconds(aux.Now)
	sd.StartTime = timeFromUnixMilliseconds(aux.StartTime)
	sd.GameOver = aux.GameOver == 1
	sd.Admin = aux.Admin == 1
	sd.TurnBased = aux.TurnBased == 1

	sd.StartTimeRaw = aux.StartTime
	return nil
}

func timeFromUnixMilliseconds(ms int64) time.Time {
	seconds := ms / 1000
	nanoseconds := (ms % 1000) * 1_000_000
	return time.Unix(seconds, nanoseconds)
}
