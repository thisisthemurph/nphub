package model

import (
	"sort"

	"github.com/labstack/gommon/log"
)

type PlayerList map[int]Player

// Get returns the player with the give UID.
// Also returns a boolean ok value indicating if the player exists.
func (pl PlayerList) Get(uid int) (Player, bool) {
	p, ok := pl[uid]
	return p, ok
}

// GetCurrent returns the player that owns the snapshot data, based on if research data is available for the player.
// Returns the first player that has research data available.
func (pl PlayerList) GetCurrent() (Player, bool) {
	for _, p := range pl {
		if p.Researching != "" {
			return p, true
		}
	}
	log.Warn("No current player found")
	return Player{}, false
}

// Sorted returns a slice of Player sorted by TotalStars descending.
func (pl PlayerList) Sorted() []Player {
	var players []Player
	for _, v := range pl {
		players = append(players, v)
	}

	sort.Slice(players, func(i, j int) bool {
		return players[i].TotalStars > players[j].TotalStars
	})

	return players
}
