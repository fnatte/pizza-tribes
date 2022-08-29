package game

import (
	"math"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/stat/distuv"
)

type Heist struct {
	Guards      int32
	Thieves     int32
	TargetCoins int32

	ThiefEvadeBonus      float64
	ThiefCapacityBonus   float64
	GuardEfficiencyBonus float64
	GuardAwarenessBonus float64
}

type HeistOutcome struct {
	Loot              int64
	SuccessfulThieves int32
	CaughtThieves     int32
	SleepingGuards    int32
}

func CalculateHeist(h Heist, rsrc rand.Source) HeistOutcome {
	guards := h.Guards
	thieves := h.Thieves

	guardsf := float64(guards)
	dist := distuv.Binomial{
		N:   guardsf,
		P:   0.075 * (1 - h.GuardAwarenessBonus),
		Src: rsrc,
	}
	sleepingGuards := MinInt32(MinInt32(int32(dist.Rand()), (guards+1)/3), thieves)

	guards = guards - sleepingGuards
	guardsf = float64(guards)
	thievesf := float64(thieves)

	guardsp := guardsf * (h.GuardEfficiencyBonus + 1) / 3
	thievesp := thievesf * (h.ThiefEvadeBonus + 1)

	dist = distuv.Binomial{
		N:   thievesf,
		P:   thievesp / (thievesp + guardsp),
		Src: rsrc,
	}
	successfulThieves := int32(dist.Rand())
	caughtThieves := MinInt32(thieves - successfulThieves, guards)
	guardsProtectingLoot := float64(MaxInt32(guards-caughtThieves, 0))
	thiefEfficiency := 0.5 + 0.5/(1+math.Pow(guardsProtectingLoot/12, 0.7))

	maxLoot := int32(float64(successfulThieves) * ThiefCapacity * (h.ThiefCapacityBonus + 1) * thiefEfficiency)
	loot := int64(MinInt32(maxLoot, h.TargetCoins))

	return HeistOutcome{
		Loot:              loot,
		SuccessfulThieves: successfulThieves,
		CaughtThieves:     caughtThieves,
		SleepingGuards:    sleepingGuards,
	}
}
