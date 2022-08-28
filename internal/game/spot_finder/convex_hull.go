package spot_finder

func convexHull(points []Vec2) []int {
	if (len(points) < 3) {
		return make([]int, 0);
	}

	// find leftmost
	lp := points[0]
	l := 0
	for i, p := range(points) {
		if p.X < lp.X {
			lp = p
			l = i
		}
	}

	hull := make([]int, 0)
	p := l;
	q := 0;
	a := 0;
	for {
		hull = append(hull, p)
		q = (p + 1) % len(points)

	for i := 0; i < len(points); i++ {
			if orientation(points[p], points[i], points[q]) < 0 {
				q = i
			}
		}

		p = q

		if p == l {
			break
		}

		a++
	}

	return hull
}
