package util

import "time"

// BoolToInt returns 1 if b is true else 0.
func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// IntToBool returns true if i is 1 and returns false for every other number.
func IntToBool(i int) bool {
	return i == 1
}

// TimeFromUnixMilliseconds returns a time.Time based on the given milliseconds.
func TimeFromUnixMilliseconds(ms int64) time.Time {
	seconds := ms / 1000
	nanoseconds := (ms % 1000) * 1_000_000
	return time.Unix(seconds, nanoseconds)
}
