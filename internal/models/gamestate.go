package models

func (gs *GameState) HasDiscovery(d ResearchDiscovery) bool {
	for _, x := range gs.Discoveries {
		if x == d {
			return true
		}
	}

	return false
}

func (gs *GameStatePatch) HasDiscovery(d ResearchDiscovery) bool {
	for _, x := range gs.Discoveries {
		if x == d {
			return true
		}
	}

	return false
}
