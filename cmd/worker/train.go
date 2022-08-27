package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fnatte/pizza-tribes/internal/game"
	"github.com/fnatte/pizza-tribes/internal/game/gamestate"
	"github.com/fnatte/pizza-tribes/internal/game/models"
	"github.com/rs/zerolog/log"
)

func (h *handler) handleTrain(ctx context.Context, userId string, m *models.ClientMessage_Train) error {
	log.Info().
		Str("userId", userId).
		Interface("Education", m.Education).
		Int32("Amount", m.Amount).
		Msg("Received train message")

	tx, err := h.updater.PerformUpdate(ctx, userId, func(gs *models.GameState, tx *gamestate.GameTx) error {
		if m.Amount <= 0 {
			return errors.New("Amount must be greater than 0")
		}

		if game.CountUneducated(gs) < m.Amount {
			return errors.New("Too few uneducated")
		}

		// Ids of mice that will go to school
		miceIds := []string{}
		for id, mouse := range gs.Mice {
			if !mouse.IsEducated && !mouse.IsBeingEducated {
				miceIds = append(miceIds, id)

				if len(miceIds) >= int(m.Amount) {
					break
				}
			}
		}

		if len(miceIds) != int(m.Amount) {
			return errors.New("Too few uneducated")
		}

		eduInfo := game.FullGameData.Educations[int32(m.Education)]
		trainTime := int64(eduInfo.TrainTime)
		cost := eduInfo.Cost * m.Amount

		if gs.Resources.Coins < cost {
			return errors.New("Not enough coins")
		}

		training := &models.Training{
			CompleteAt: time.Now().UnixNano() + trainTime*1e9,
			Education:  m.Education,
			Amount:     m.Amount,
		}

		u := tx.Users[userId]
		for _, id := range miceIds {
			u.SetMouseIsBeingEducated(id, true)
		}
		u.SetCoins(gs.Resources.Coins - cost)
		u.AppendTrainingQueue(training)

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to perform train update: %w", err)
	}

	err = h.sendGameTx(ctx, tx)
	if err != nil {
		return fmt.Errorf("failed to send game tx: %w", err)
	}

	return nil
}
