package model

import (
	"encoding/json"
	"strconv"
)

type Player struct {
	UID             int                     `json:"uid"`
	AI              bool                    `json:"ai"`
	Alias           string                  `json:"alias"`
	AvatarID        int                     `json:"avatar"`
	HomeStarUID     int                     `json:"huid"`
	Cash            int                     `json:"cash"` // Only available for the current player
	Tech            map[string]TechLevel    `json:"tech"`
	Researching     string                  `json:"researching"`      // Only available for the current player
	ResearchingNext string                  `json:"researching_next"` // Only available for the current player
	TotalIndustry   int                     `json:"total_industry"`
	TotalEconomy    int                     `json:"total_economy"`
	TotalScience    int                     `json:"total_science"`
	TotalStars      int                     `json:"total_stars"`
	TotalFleets     int                     `json:"total_fleets"`
	TotalStrength   int                     `json:"total_strength"` // Total ships the player has
	War             map[int]PlayerWarStatus `json:"war"`
	CountdownToWar  map[int]int             `json:"countdown_to_war"` // An object containing all player IDs and the number of ticks until war starts, if a permanent alliance has ended
	Ready           bool                    `json:"ready"`
	Regard          int                     `json:"regard"` // The AI’s opinion of the player. Note that this may be present for non-AI players.
	Conceded        bool                    `json:"conceded"`
	StarsAbandoned  int                     `json:"stars_abandoned"` // Number of stars abandoned this production round (note: can’t be higher than 1, resets to 0 at prod)
	MissedTurns     int                     `json:"missed_turns"`
	KarmaToGive     int                     `json:"karma_to_give"`
}

func (p *Player) UnmarshalJSON(data []byte) error {
	type Alias Player
	aux := struct {
		AI             int `json:"ai"`
		Conceded       int `json:"conceded"`
		Ready          int `json:"ready"`
		War            map[string]int
		CountdownToWar map[string]int
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	for uid, status := range aux.War {
		playerUID, err := strconv.Atoi(uid)
		if err != nil {
			return err
		}
		p.War[playerUID] = PlayerWarStatus(status)
	}

	for uid, ticks := range aux.CountdownToWar {
		playerUID, err := strconv.Atoi(uid)
		if err != nil {
			return err
		}
		p.CountdownToWar[playerUID] = ticks
	}

	p.AI = aux.AI == 1
	p.Conceded = aux.Conceded == 1
	p.Ready = aux.Ready == 1
	return nil
}

type PlayerWarStatus int

func (pws PlayerWarStatus) String() string {
	switch pws {
	case 0:
		return "Formal Alliance"
	case 1:
		return "Alliance Offered"
	case 2:
		return "Alliance Offered (paid)"
	case 3:
		return "At War"
	default:
		return "Unknown War Status"
	}
}
