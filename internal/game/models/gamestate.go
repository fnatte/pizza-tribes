package models

import "github.com/rs/xid"

func (gs *GameState) HasDiscovery(d ResearchDiscovery) bool {
	for _, x := range gs.Discoveries {
		if x == d {
			return true
		}
	}

	return false
}

func NewGameState() *GameState {
	return &GameState{
		Resources:  &GameState_Resources{},
		Lots: map[string]*GameState_Lot{
			"4": {
				Building: Building_TOWN_CENTRE,
			},
		},
		Discoveries: []ResearchDiscovery{},
		Mice:        map[string]*Mouse{},
		Quests: map[string]*QuestState{
			"1": {},
		},
		PizzaPrice: 1,
	}
}

func (gsp *GameStatePatch) ToServerMessage() *ServerMessage {
	return &ServerMessage{
		Id: xid.New().String(),
		Payload: &ServerMessage_StateChange{
			StateChange: &GameStatePatch{
				GameState: gsp.GameState,
				PatchMask: gsp.PatchMask,
			},
		},
	}
}

func (gs *GameState) GetValidPizzaPrice() int32 {
	p := gs.PizzaPrice

	if p < 1 {
		return 1
	}
	if p > 15 {
		return 15
	}

	return p
}

func (gs *GameState) CountTownPopulationEducations() map[Education]int32 {
	res := map[Education]int32{}
	for _, m := range gs.Mice {
		if m.IsEducated {
			res[m.Education] = res[m.Education] + 1
		}
	}

	for _, t := range gs.TravelQueue {
		if t.Thieves > 0 {
			res[Education_THIEF] = res[Education_THIEF] - t.Thieves
		}
	}

	return res
}

