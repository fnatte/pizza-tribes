package models

import (
	"github.com/rs/xid"
)

// Deprecated: Use ToServerMessage instead.
func (gs *GameState) ToStateChangeMessage() *ServerMessage {
	return gs.ToServerMessage()
}

func (gs *GameState) ToServerMessage() *ServerMessage {
	return &ServerMessage{
		Id: xid.New().String(),
		Payload: &ServerMessage_StateChange{
			StateChange: &GameStatePatch{
				GameState: gs,
			},
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
