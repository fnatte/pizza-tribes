package internal

import (
	"math"
	"time"
)

// Calculates the arrival time from now.
// Returns the travel time in nanoseconds.
func CalculateArrivalTime(fromX, fromY, toX, toY int32, speed time.Duration) int64 {
	base := 3.0
	dx := toX - fromX
	dy := toY - fromY
	distance := math.Sqrt(float64(dx*dx) + float64(dy*dy)) + base
	travelTime := distance * speed.Seconds()
	return time.Now().UnixNano() + int64(travelTime*1e9)
}
