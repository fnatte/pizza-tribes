package spot_finder

import (
	"math"
	"math/rand"
)

type Point struct {
	X int
	Y int
}

const m = 3

func randomOffset(v Vec2) Vec2 {
	return NewVec2(v.X+(rand.Float64()-0.5)*3, v.Y+(rand.Float64()-0.5)*3)
}

var firstThree = []Vec2{
	{X: m, Y: 0},
	{X: 0, Y: m},
	{X: m, Y: m},
}

func findFirstThree(existing []Vec2) Point {
	available := make([]Vec2, len(firstThree))
	copy(available, firstThree)
	for i := len(available) - 1; i >= 0; i-- {
		a := available[i]
		for _, b := range existing {
			if int(a.X) == int(b.X) && int(a.Y) == int(b.Y) {
				available = append(available[:i], available[i+1:]...)
				break
			}
		}
	}

	return Point{ X: int(available[0].X), Y: int(available[0].Y)}
}

func findNextFromHull(existing []Vec2) Point {
	hull := convexHull(existing)
	hull = append(hull, hull[0])
	var p Vec2
	var d = math.MaxFloat64
	for i := 0; i < len(hull)-1; i++ {
		p1 := existing[hull[i]]
		p2 := existing[hull[i+1]]
		px := (p2.X-p1.X)/2 + p1.X
		py := (p2.Y-p1.Y)/2 + p1.Y
		nd := px*px + py*py
		if nd < d {
			puv := unit(perp(NewVec2(p2.X-p1.X, p2.Y-p1.Y)))
			p = randomOffset(NewVec2(px+puv.X*3, py+puv.Y*3))
			d = nd
		}
	}

	return Point{ X: int(p.X), Y: int(p.Y) }
}

func FindSpotForNewTown(existing []Vec2, blocked []Point) Point {
	if len(existing) < 3 {
		return findFirstThree(existing)
	}

	v2s := make([]Vec2, len(existing))
	copy(v2s, existing)

	for i := 0; i < 10_000_000; i++ {
		p := findNextFromHull(v2s)

		isBlocked := false
		for _, p2 := range blocked {
			if p.X == p2.X && p.Y == p2.Y {
				isBlocked = true
				break
			}
		}

		if !isBlocked {
			return p
		}

		v2s = append(v2s, Vec2{ X: float64(p.X), Y: float64(p.Y) })
	}

	panic("ended up in infinte loop when finding spot for new town")
}

