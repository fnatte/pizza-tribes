package main

import (
	"context"
	"fmt"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/fnatte/pizza-tribes/internal/gamestate"
	"github.com/fnatte/pizza-tribes/internal/models"
	"github.com/fnatte/pizza-tribes/internal/protojson"
	"github.com/rs/xid"
	"github.com/rs/zerolog/log"
)

type handler struct {
	rdb   internal.RedisClient
	world *internal.WorldService
}

func (h *handler) Handle(ctx context.Context, senderId string, m *models.ClientMessage) {
	var err error

	state, err := h.world.GetState(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to handle message")
		return
	}

	if _, ok := state.Type.(*models.WorldState_Started_); !ok {
		log.Info().Msg("Announcing world state.")
		h.send(ctx, senderId, &models.ServerMessage{
			Id:      xid.New().String(),
			Payload: &models.ServerMessage_WorldState{
				WorldState: state,
			},
		})
		return
	}

	switch x := m.Type.(type) {
	case *models.ClientMessage_Tap_:
		err = h.handleTap(ctx, senderId, x.Tap)
	case *models.ClientMessage_ConstructBuilding_:
		err = h.handleConstructBuilding(ctx, senderId, x.ConstructBuilding)
	case *models.ClientMessage_Train_:
		err = h.handleTrain(ctx, senderId, x.Train)
	case *models.ClientMessage_Steal_:
		err = h.handleSteal(ctx, senderId, x.Steal)
	case *models.ClientMessage_ReadReport_:
		err = h.handleReadReport(ctx, senderId, x.ReadReport)
	case *models.ClientMessage_UpgradeBuilding_:
		err = h.handleUpgrade(ctx, senderId, x.UpgradeBuilding)
	case *models.ClientMessage_RazeBuilding_:
		err = h.handleRazeBuilding(ctx, senderId, x.RazeBuilding)
	case *models.ClientMessage_CancelRazeBuilding_:
		err = h.handleCancelRazeBuilding(ctx, senderId, x.CancelRazeBuilding)
	case *models.ClientMessage_StartResearch_:
		err = h.handleStartResearch(ctx, senderId, x.StartResearch)
	case *models.ClientMessage_ReschoolMouse_:
		err = h.handleReschoolMouse(ctx, senderId, x.ReschoolMouse)
	case *models.ClientMessage_RenameMouse_:
		err = h.handleRenameMouse(ctx, senderId, x.RenameMouse)
	case *models.ClientMessage_OpenQuest_:
		err = h.handleOpenQuest(ctx, senderId, x.OpenQuest)
	case *models.ClientMessage_ClaimQuestReward_:
		err = h.handleClaimQuestReward(ctx, senderId, x.ClaimQuestReward)
	case *models.ClientMessage_CompleteVisitHelpPageQuest_:
		err = h.handleCompleteVisitHelpPageQuest(ctx, senderId, x.CompleteVisitHelpPageQuest)
	default:
		log.Info().Str("senderId", senderId).Msg("Received message")
	}

	if err != nil {
		log.Error().Err(err).Msg("failed to handle message")
	}
}

func (h *handler) fetchAndUpdateTimestamp(ctx context.Context, userId string) (int64, error) {
	s, err := h.rdb.JsonGet(ctx, fmt.Sprintf("user:%s:gamestate", userId), ".").Result()
	if err != nil {
		return 0, err
	}

	gs := models.GameState{}
	if err = protojson.Unmarshal([]byte(s), &gs); err != nil {
		return 0, err
	}

	return internal.SetNextUpdate(h.rdb, ctx, userId, &gs)
}

func (h *handler) sendFullStateUpdate(ctx context.Context, senderId string) {
	s, err := h.rdb.JsonGet(ctx, fmt.Sprintf("user:%s:gamestate", senderId), ".").Result()
	if err != nil {
		log.Error().Err(err).Msg("Failed to send full state update")
		return
	}

	gs := models.GameState{}
	if err = protojson.Unmarshal([]byte(s), &gs); err != nil {
		log.Error().Err(err).Msg("Failed to send full state update")
		return
	}

	msg := gs.ToStateChangeMessage()
	err = h.send(ctx, senderId, msg)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send state update")
		return
	}

	msg = internal.CalculateStats(&gs).ToServerMessage()
	err = h.send(ctx, senderId, msg)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send stats message")
		return
	}
}

func (h *handler) sendGameTx(ctx context.Context, tx *gamestate.GameTx) error {
	errs := []error{}

	for uid, u := range tx.Users {
		/*
		jsonPatch := []*models.JsonPatchOp{}
		for _, op := range u.GsPatch.Ops {
			val, err := json.Marshal(op.Value)
			if err != nil {
				return err
			}

			jsonPatch = append(jsonPatch, &models.JsonPatchOp{
				From:  op.From,
				Op:    op.Op,
				Path:  op.Path,
				Value: string(val),
			})
		}
*/

		err := h.send(ctx, uid, &models.ServerMessage{
			Id: xid.New().String(),
			Payload: &models.ServerMessage_StateChange2{
				StateChange2: &models.ServerMessage_GameStatePatch2{
					JsonPatch: u.GsPatch.Ops,
				},
			},
		})
		if err != nil {
			errs = append(errs, err)
		}

		if u.StatsInvalidated {
			msg := internal.CalculateStats(u.Gs).ToServerMessage()
			err = h.send(ctx, uid, msg)
			if err != nil {
				errs = append(errs, err)
			}
		}

		if u.NextUpdateInvalidated {
			_, err = internal.SetNextUpdate(h.rdb, ctx, uid, u.Gs)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors when sending game tx")
	}

	return nil
}

func (h *handler) send(ctx context.Context, senderId string, m *models.ServerMessage) error {
	b, err := protojson.Marshal(m)
	if err != nil {
		return err
	}

	h.rdb.RPush(ctx, "wsout", &internal.OutgoingMessage{
		ReceiverId: senderId,
		Body:       string(b),
	})

	return nil
}
