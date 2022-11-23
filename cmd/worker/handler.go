package main

import (
	"context"
	"fmt"
	"time"

	"github.com/fnatte/pizza-tribes/internal/game"
	"github.com/fnatte/pizza-tribes/internal/game/gamestate"
	"github.com/fnatte/pizza-tribes/internal/game/models"
	"github.com/fnatte/pizza-tribes/internal/game/persist"
	"github.com/fnatte/pizza-tribes/internal/game/protojson"
	"github.com/fnatte/pizza-tribes/internal/game/redis"
	"github.com/rs/xid"
	"github.com/rs/zerolog/log"
)

type handler struct {
	rdb         redis.RedisClient
	world       *game.WorldService
	gsRepo      persist.GameStateRepository
	reportsRepo persist.ReportsRepository
	userRepo    persist.GameUserRepository
	marketRepo  persist.MarketRepository
	updater     gamestate.Updater
	speed       float64
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
			Id: xid.New().String(),
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
	case *models.ClientMessage_CompleteQuest_:
		err = h.handleCompleteQuest(ctx, senderId, x.CompleteQuest)
	case *models.ClientMessage_SaveMouseAppearance_:
		err = h.handleSaveMouseAppearance(ctx, senderId, x.SaveMouseAppearance)
	case *models.ClientMessage_SetAmbassadorMouse_:
		err = h.handleSetAmbassadorMouse(ctx, senderId, x.SetAmbassadorMouse)
	case *models.ClientMessage_SetPizzaPrice_:
		err = h.handleSetPizzaPrice(ctx, senderId, x.SetPizzaPrice)
	case *models.ClientMessage_BuyGeniusFlash_:
		err = h.handleBuyGeniusFlash(ctx, senderId, x.BuyGeniusFlash)
	default:
		log.Debug().Str("senderId", senderId).Msg("Received message")
	}

	if err != nil {
		log.Error().Err(err).Msg("failed to handle message")
	}

	// Update users latest activity
	err = h.userRepo.SetUserLatestActivity(ctx, senderId, time.Now().UnixNano())
	if err != nil {
		log.Error().Err(err).Msg("failed to update user's latest activity")
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

	return game.SetNextUpdate(h.rdb, ctx, userId, &gs)
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

	err = h.send(ctx, senderId, &models.ServerMessage{
		Id: xid.New().String(),
		Payload: &models.ServerMessage_StateChange{
			StateChange: &models.GameStatePatch{
				GameState: &gs,
			},
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to send state update")
		return
	}

	h.sendStats(ctx, senderId, &gs)
}

func (h *handler) sendStats(ctx context.Context, userId string, gs *models.GameState) error {
	worldState, err := h.world.GetState(ctx)
	if err != nil {
		return fmt.Errorf("failed to send full state update: %w", err)
	}

	globalDemandScore, err := h.marketRepo.GetGlobalDemandScore(ctx)
	if err != nil {
		return fmt.Errorf("failed to send full state update: %w", err)
	}

	userCount, err := h.userRepo.GetUserCount(ctx)
	if err != nil {
		return fmt.Errorf("failed to send full state update: %w", err)
	}

	msg := game.CalculateStats(gs, globalDemandScore, worldState, userCount).ToServerMessage()
	err = h.send(ctx, userId, msg)
	if err != nil {
		return fmt.Errorf("failed to send full state update: %w", err)
	}

	return nil
}

func (h *handler) sendGameTx(ctx context.Context, tx *gamestate.GameTx) error {
	errs := []error{}

	for uid, u := range tx.Users {
		err := h.send(ctx, uid, &models.ServerMessage{
			Id: xid.New().String(),
			Payload: &models.ServerMessage_StateChange{
				StateChange: &models.GameStatePatch{
					GameState: u.Gs,
					PatchMask: &models.PatchMask{
						Paths: u.PatchMask.Paths,
					},
				},
			},
		})
		if err != nil {
			errs = append(errs, err)
		}

		if u.StatsInvalidated {
			err = h.sendStats(ctx, uid, u.Gs)
			if err != nil {
				errs = append(errs, err)
			}
		}

		if u.NextUpdateInvalidated {
			_, err = game.SetNextUpdate(h.rdb, ctx, uid, u.Gs)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	if len(errs) > 0 {
		if len(errs) == 1 {
			return fmt.Errorf("error when sending game tx: %w", errs[0])
		}
		return fmt.Errorf("errors when sending game tx")
	}

	return nil
}

func (h *handler) send(ctx context.Context, senderId string, m *models.ServerMessage) error {
	b, err := protojson.Marshal(m)
	if err != nil {
		return err
	}

	h.rdb.RPush(ctx, "wsout", &game.OutgoingMessage{
		ReceiverId: senderId,
		Body:       string(b),
	})

	return nil
}
