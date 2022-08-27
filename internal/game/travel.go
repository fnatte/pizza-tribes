package game

import (
	"math"
	"time"

	"github.com/fnatte/pizza-tribes/internal/game/models"
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

func GetThiefSpeed(gs *models.GameState) time.Duration {
	nsec := float64(int64(ThiefSpeed))

	if gs.HasDiscovery(models.ResearchDiscovery_BOOTS_OF_HASTE) {
		nsec *= 1.25
	}

	return time.Duration(int64(nsec))
}
