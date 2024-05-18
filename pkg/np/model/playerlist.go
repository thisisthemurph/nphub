package model

import "sort"

type PlayerList map[int]Player

// Get returns the player with the give UID.
// Also returns a boolean ok value indicating if the player exists.
func (pl PlayerList) Get(uid int) (Player, bool) {
	p, ok := pl[uid]
	return p, ok
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
