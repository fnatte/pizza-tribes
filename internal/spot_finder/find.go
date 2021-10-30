package spot_finder

import (
	"math"
	"math/rand"
)

const m = 3

func randomOffset(v Vec2) Vec2 {
	return NewVec2(v.X + (rand.Float64() - 0.5) * 3, v.Y + (rand.Float64() - 0.5) * 3)
}

func FindSpotForNewTown(existing []Vec2) (x, y int) {
	if len(existing) < 3 {
		switch len(existing) {
		case 0: return m, 0
		case 1: return 0, m
		case 2: return m, m
		}
	}

	hull := convexHull(existing)
	hull = append(hull, hull[0])
	var p Vec2
	var d = math.MaxFloat64
	for i := 0; i < len(hull) - 1; i++ {
		p1 := existing[hull[i]]
		p2 := existing[hull[i + 1]]
		px := (p2.X - p1.X) / 2 + p1.X
		py := (p2.Y - p1.Y) / 2 + p1.Y
		nd := px * px + py * py
		if nd < d {
			puv := unit(perp(NewVec2(p2.X - p1.X, p2.Y - p1.Y)))
			p = randomOffset(NewVec2(px + puv.X * 3, py + puv.Y * 3))
			d = nd
		}
	}

	return int(p.X), int(p.Y)
}

