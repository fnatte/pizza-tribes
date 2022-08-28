package spot_finder

import (
	"math"
)

type Vec2 struct {
	X float64
	Y float64
}

func NewVec2(x, y float64) Vec2 {
	return Vec2{X:x, Y: y}
}

func orientation(a, b, c Vec2) float64 {
	return (b.Y - a.Y) * (c.X - b.X) - (b.X - a.X) * (c.Y - b.Y)
}

func perp(v Vec2) Vec2 {
	return NewVec2(v.Y, -v.X)
}

func unit(v Vec2) Vec2 {
	d := math.Sqrt(v.X * v.X + v.Y * v.Y)
	return NewVec2(v.X / d, v.Y / d)
}


