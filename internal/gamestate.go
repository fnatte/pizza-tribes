package internal

import "time"

func (gs *GameState) LoadProtoJson(b []byte) error {
	return protojsonu.Unmarshal(b, gs)
}

func (gs *GameState) GetCompletedTravels() (res []*Travel) {
	now := time.Now().UnixNano()

	for _, t := range gs.TravelQueue {
		if t.ArrivalAt > now {
			break
		}

		res = append(res, t)
	}

	return res
}

