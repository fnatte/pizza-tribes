package internal

import (
	"time"

	"google.golang.org/protobuf/encoding/protojson"
)

func NewInt64(i int64) *int64 { return &i }
func NewString(s string) *string { return &s }

var protojsonu = protojson.UnmarshalOptions{
	DiscardUnknown: true,
}

func CountBuildings(gs *GameState) (counts map[int32]int32) {
	counts = map[int32]int32{}
	for _, lot := range gs.Lots {
		counts[int32(lot.Building)] = counts[int32(lot.Building)] + 1
	}
	return counts
}

func CountBuildingsUnderConstruction(gs *GameState) (counts map[int32]int32) {
	counts = map[int32]int32{}
	for _, c := range gs.ConstructionQueue {
		counts[int32(c.Building)] = counts[int32(c.Building)] + 1
	}
	return counts
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func GetNextUpdateTimestamp (gs *GameState) int64 {
	t := time.Now().Add(10 * time.Second).UnixNano()
	if len(gs.ConstructionQueue) > 0 {
		t = min(t, gs.ConstructionQueue[0].CompleteAt)
	}
	if len(gs.TrainingQueue) > 0 {
		t = min(t, gs.TrainingQueue[0].CompleteAt)
	}
	return t
}
