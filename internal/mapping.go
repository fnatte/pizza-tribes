package internal

import "github.com/rs/xid"

func (gs *GameState) ToStateChangeMessage() *ServerMessage {
	lotsPatch := map[string]*GameStatePatch_LotPatch{}
	for lotId, lot := range gs.Lots {
		lotsPatch[lotId] = &GameStatePatch_LotPatch{
			Building: &lot.Building,
		}
	}

	p := &GameStatePatch{
		Lots: lotsPatch,
		Resources: &GameStatePatch_ResourcesPatch{
			Coins:  NewInt64(gs.Resources.Coins),
			Pizzas: NewInt64(gs.Resources.Pizzas),
		},
		Population: &GameStatePatch_PopulationPatch{
			Unemployed: NewInt64(gs.Population.Unemployed),
			Chefs:      NewInt64(gs.Population.Chefs),
			Salesmice:  NewInt64(gs.Population.Salesmice),
			Guards:     NewInt64(gs.Population.Guards),
			Thieves:    NewInt64(gs.Population.Thieves),
		},
		TrainingQueue: gs.TrainingQueue,
		TrainingQueuePatched: true,
	}

	return &ServerMessage{
		Id: xid.New().String(),
		Payload: &ServerMessage_StateChange{
			StateChange: p,
		},
	}
}
