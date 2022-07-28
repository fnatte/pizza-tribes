package models

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
		Population: &GameState_Population{},
		Resources:  &GameState_Resources{},
		Lots: map[string]*GameState_Lot{
			"2": {
				Building: Building_TOWN_CENTRE,
			},
		},
		Discoveries: []ResearchDiscovery{},
		Mice:        map[string]*Mouse{},
		Quests: map[string]*QuestState{
			"1": {},
		},
	}
}
