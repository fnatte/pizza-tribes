package internal

import (
	"github.com/rs/xid"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

func (gs *GameState) ToStateChangeMessage() *ServerMessage {
	lotsPatch := map[string]*GameStatePatch_LotPatch{}
	for lotId, lot := range gs.Lots {
		lotsPatch[lotId] = &GameStatePatch_LotPatch{
			Building: lot.Building,
		}
	}

	pop := &GameStatePatch_PopulationPatch{}

	if gs.Population != nil {
		pop.Uneducated = &wrapperspb.Int64Value{ Value: gs.Population.Uneducated }
		pop.Chefs =      &wrapperspb.Int64Value{ Value: gs.Population.Chefs }
		pop.Salesmice =  &wrapperspb.Int64Value{ Value: gs.Population.Salesmice }
		pop.Guards =     &wrapperspb.Int64Value{ Value: gs.Population.Guards }
		pop.Thieves =    &wrapperspb.Int64Value{ Value: gs.Population.Thieves }
	}

	p := &GameStatePatch{
		Lots: lotsPatch,
		Resources: &GameStatePatch_ResourcesPatch{
			Coins: &wrapperspb.Int64Value{ Value: gs.Resources.Coins },
			Pizzas: &wrapperspb.Int64Value{ Value: gs.Resources.Pizzas },
		},
		Population: pop,
		TrainingQueue: gs.TrainingQueue,
		TrainingQueuePatched: true,
		ConstructionQueue: gs.ConstructionQueue,
		ConstructionQueuePatched: true,
	}

	return &ServerMessage{
		Id: xid.New().String(),
		Payload: &ServerMessage_StateChange{
			StateChange: p,
		},
	}
}
