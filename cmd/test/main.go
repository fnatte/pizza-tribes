package main

import (
	"fmt"
	"time"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)


func main() {
	guards := float64(50)
	thieves := float64(200)
	// maxLoot := thieves * 3_000
	// loot := int64(internal.MinInt32(maxLoot, gsTarget.Resources.Coins))
	dist := distuv.Binomial{
		N: thieves,
		P: thieves / (thieves + guards),
		Src: rand.NewSource(uint64(time.Now().UnixNano())),
	}
	successfulThieves := dist.Rand()

	fmt.Printf("%f\n", successfulThieves)
}
