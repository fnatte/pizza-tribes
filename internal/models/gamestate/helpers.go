package gamestate

func (gs *GameState) HasDiscovery(d ResearchDiscovery) bool {
	for _, x := range gs.Discoveries {
		if x == d {
			return true
		}
	}

	return false
}
