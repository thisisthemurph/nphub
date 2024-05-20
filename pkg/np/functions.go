package np

import (
	"time"
)

// CalculateNextTickTime calculates the time of the next tick for the given game params.
func CalculateNextTickTime(gameStartTime time.Time, gameTickRate int) (time.Time, error) {
	gameDurationMins := int64(time.Since(gameStartTime).Minutes())
	elapsedTicks := gameDurationMins / int64(gameTickRate)
	lastTickTime := gameStartTime.Add(time.Duration(elapsedTicks*int64(gameTickRate)) * time.Minute)

	return lastTickTime.Add(time.Duration(gameTickRate) * time.Minute), nil
}
