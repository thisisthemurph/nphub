package model

import "encoding/json"

type ScanningData struct {
	Fleets            map[int]Fleet  `json:"fleets"`
	FleetSpeed        float32        `json:"fleet_speed"`
	Paused            bool           `json:"paused"`
	Productions       int            `json:"productions"`
	TickFragment      float32        `json:"tick_fragment"`
	Now               int64          `json:"now"`
	TickRate          int            `json:"tick_rate"`
	ProductionRate    int            `json:"production_rate"`
	Stars             map[int]Star   `json:"stars"`
	StarsForVictory   int            `json:"stars_for_victory"`
	GameOver          bool           `json:"game_over"` // int in original JSON
	Started           bool           `json:"started"`
	StartTime         int64          `json:"start_time"`
	TotalStars        int            `json:"total_stars"`
	ProductionCounter int            `json:"production_counter"`
	TradeScanned      int            `json:"trade_scanned"`
	Tick              int            `json:"tick"`
	TradeCost         int            `json:"trade_cost"`
	Name              string         `json:"name"`
	PlayerUID         int            `json:"player_uid"`
	Admin             bool           `json:"admin"`      // int in original JSON
	TurnBased         bool           `json:"turn_based"` // int in original JSON
	War               bool           `json:"war"`        // int in original JSON
	Players           map[int]Player `json:"players"`
	TurnBasedTimeOut  bool           `json:"turn_based_time_out"` // int in original JSON
}

func NewScanningData(data []byte) (ScanningData, error) {
	var scd ScanningData
	err := json.Unmarshal(data, &scd)
	return scd, err
}

type Fleet struct {
	// TODO: Complete Fleet struct
	UID  int     `json:"uid"`
	Name string  `json:"n"`
	X    float32 `json:"x"`
	Y    float32 `json:"y"`
}

type Star struct {
	// TODO: Complete Star struct
	UID  int     `json:"uid"`
	Name string  `json:"n"`
	PuID int     `json:"puid"`
	V    int     `json:"v"`
	X    float32 `json:"x"`
	Y    float32 `json:"y"`
}

type Player struct {
	TotalIndustry int                  `json:"total_industry"`
	Regard        int                  `json:"regard"`
	TotalScience  int                  `json:"total_science"`
	UID           int                  `json:"uid"`
	AI            bool                 `json:"ai"` // int in original JSON
	HuID          int                  `json:"huid"`
	TotalStars    int                  `json:"total_stars"`
	TotalFleets   int                  `json:"total_fleets"`
	TotalStrength int                  `json:"total_strength"`
	Alias         string               `json:"alias"`
	Tech          map[string]TechLevel `json:"tech"`
	AvatarID      int                  `json:"avatar"`
	Conceded      int                  `json:"conceded"`
	Ready         bool                 `json:"ready"` // int in original JSON
	TotalEconomy  int                  `json:"total_economy"`
	MissedTurns   int                  `json:"missed_turns"`
	KarmaToGive   int                  `json:"karma_to_give"`
}

type TechLevel struct {
	Value float32 `json:"value"`
	Level int     `json:"level"`
}
