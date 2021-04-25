package internal

import "github.com/rs/xid"

func (gs *GameState) ToStateChangeMessage() *ServerMessage {
	lotsPatch := map[string]*GameStatePatch_LotPatch{}
	for lotId, lot := range(gs.Lots) {
		lotsPatch[lotId] = &GameStatePatch_LotPatch{
			Building: &lot.Building,
		}
	}

	return &ServerMessage{
		Id: xid.New().String(),
		Payload: &ServerMessage_StateChange{
			StateChange: &GameStatePatch{
				Lots: lotsPatch,
				Resources: &GameStatePatch_ResourcesPatch{
					Coins:  NewInt64(gs.Resources.Coins),
					Pizzas: NewInt64(gs.Resources.Pizzas),
				},
			},
		},
	}
}
