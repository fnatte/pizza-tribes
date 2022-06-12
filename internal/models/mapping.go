package models

import (
	"github.com/rs/xid"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

// Deprecated: Use ToServerMessage instead.
func (gs *GameState) ToStateChangeMessage() *ServerMessage {
	return gs.ToServerMessage()
}

func (gs *GameState) ToServerMessage() *ServerMessage {
	lotsPatch := map[string]*GameStatePatch_LotPatch{}
	for lotId, lot := range gs.Lots {
		lotsPatch[lotId] = &GameStatePatch_LotPatch{
			Building: lot.Building,
			TappedAt: lot.TappedAt,
			Level:    lot.Level,
			Taps:     lot.Taps,
			Streak:   lot.Streak,
		}
	}

	pop := &GameStatePatch_PopulationPatch{}

	if gs.Population != nil {
		pop.Uneducated = &wrapperspb.Int32Value{Value: gs.Population.Uneducated}
		pop.Chefs = &wrapperspb.Int32Value{Value: gs.Population.Chefs}
		pop.Salesmice = &wrapperspb.Int32Value{Value: gs.Population.Salesmice}
		pop.Guards = &wrapperspb.Int32Value{Value: gs.Population.Guards}
		pop.Thieves = &wrapperspb.Int32Value{Value: gs.Population.Thieves}
		pop.Publicists = &wrapperspb.Int32Value{Value: gs.Population.Publicists}
	}

	mice := map[string]*GameStatePatch_MousePatch{}
	if gs.Mice != nil {
		for id, m := range gs.Mice {
			mice[id] = m.ToPatch(false)
		}
	}

	p := &GameStatePatch{
		Lots: lotsPatch,
		Resources: &GameStatePatch_ResourcesPatch{
			Coins:  &wrapperspb.Int32Value{Value: gs.Resources.Coins},
			Pizzas: &wrapperspb.Int32Value{Value: gs.Resources.Pizzas},
		},
		Population:               pop,
		TrainingQueue:            gs.TrainingQueue,
		TrainingQueuePatched:     true,
		ConstructionQueue:        gs.ConstructionQueue,
		ConstructionQueuePatched: true,
		TownX:                    &wrapperspb.Int32Value{Value: gs.TownX},
		TownY:                    &wrapperspb.Int32Value{Value: gs.TownY},
		TravelQueue:              gs.TravelQueue,
		TravelQueuePatched:       true,
		DiscoveriesPatched:       true,
		Discoveries:              gs.Discoveries,
		ResearchQueuePatched:     true,
		ResearchQueue:            gs.ResearchQueue,
		Mice:                     mice,
	}

	return &ServerMessage{
		Id: xid.New().String(),
		Payload: &ServerMessage_StateChange{
			StateChange: p,
		},
	}
}

func (stats *Stats) ToServerMessage() *ServerMessage {
	return &ServerMessage{
		Id: xid.New().String(),
		Payload: &ServerMessage_Stats{
			Stats: stats,
		},
	}
}
