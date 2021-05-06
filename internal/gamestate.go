package internal

func (gs *GameState) LoadProtoJson(b []byte) error {
	return protojsonu.Unmarshal(b, gs)
}

