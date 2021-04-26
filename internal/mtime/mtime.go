package mtime

import "time"

const (
	ScaleFactor = 24

	secondsPerMinute = 60
	secondsPerHour   = 60 * secondsPerMinute
	secondsPerDay    = 24 * secondsPerHour
)

func Now() int64 {
	return HumanToMouse(time.Now())
}

func Clock(mt int64) (hour, min, sec int) {
	sec = int(mt % secondsPerDay)
	hour = sec / secondsPerHour
	sec -= hour * secondsPerHour
	min = sec / secondsPerMinute
	sec -= min * secondsPerMinute
	return
}

func ClockNow() (hour, min, sec int) {
	return Clock(Now())
}

func HumanToMouse(ht time.Time) int64 {
	return ht.Unix() * ScaleFactor
}

func HumanUnixToMouse(ht int64) int64 {
	return ht * ScaleFactor
}

