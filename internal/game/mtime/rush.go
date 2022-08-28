package mtime

func IsRush(t int64) bool {
	h, _, _ := Clock(HumanUnixToMouse(t))
	return (h >= 11 && h < 13) || (h >= 18 && h < 21)
}

func GetRush(then, now int64) (rush, offpeak int64) {
	dt := now - then

	// If extrapolating less than a minute, just check if midtime is rush
	if dt < 60 {
		if IsRush(then + dt/2) {
			return dt, 0
		} else {
			return 0, dt
		}
	}

	// TODO: This could be done a lot smarter/cheaper.

	n := int64(150*4)
	s := float64(dt) / float64(n)
	rushf := 0.0
	offpeakf := 0.0
	for i := int64(0); i < n; i++ {
		if IsRush(then + int64(float64(i) * s)) {
			rushf = rushf + s
		} else {
			offpeakf = offpeakf + s
		}
	}

	return int64(rushf), int64(offpeakf)
}
