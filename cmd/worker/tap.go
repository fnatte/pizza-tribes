package main

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/fnatte/pizza-tribes/internal/gamestate"
	"github.com/fnatte/pizza-tribes/internal/models"
)

const MAX_TAP_STREAK = 12

func (h *handler) handleTap(ctx context.Context, userId string, m *models.ClientMessage_Tap) error {
	now := time.Now().UnixNano()

	tx, err := h.updater.PerformUpdate(ctx, userId, func(gs *models.GameState, tx *gamestate.GameTx) error {
		lot := gs.Lots[m.LotId]
		if lot == nil {
			return fmt.Errorf("invalid lot: %s", m.LotId)
		}

		// Check tap interval
		nextTapAt := lot.TappedAt + (500 * time.Millisecond).Nanoseconds()
		if nextTapAt > now {
			return fmt.Errorf("tapped to soon, next tap at %d", nextTapAt)
		}

		// Set reset time to the beginning of the next hour after TappedAt,
		// and streak time to the hour after that.
		resetTime := time.Unix(0, lot.TappedAt).Add(1 * time.Hour).Truncate(1 * time.Hour)
		resetStreakTime := resetTime.Add(1 * time.Hour)

		// Reset streak if we are past the reset streak time
		if time.Now().After(resetStreakTime) {
			lot.Streak = 0
		}

		taps := gs.Lots[m.LotId].Taps
		streak := gs.Lots[m.LotId].Streak

		// Reset taps if next hour
		if time.Now().After(resetTime) {
			taps = 0
		}

		if taps >= 10 {
			return fmt.Errorf("tap is maxed out this hour")
		}

		// Determine what resource to increase and how much
		var incrType string
		var incrAmount int32
		factor := math.Sqrt(float64(lot.Level+1) * float64(lot.Streak+1))
		switch lot.Building {
		case models.Building_KITCHEN:
			incrType = "pizzas"
			incrAmount = int32(math.Round(80*factor/5) * 5)
		case models.Building_SHOP:
			incrType = "coins"
			incrAmount = int32(math.Round(35*factor/5) * 5)
		default:
			return fmt.Errorf("this building cannot be tapped")
		}

		taps = taps + 1

		if taps == 10 {
			streak = streak + 1
		}

		u := tx.Users[userId]
		u.SetTappedAt(m.LotId, int64(now))
		u.SetTaps(m.LotId, taps)
		u.SetStreak(m.LotId, streak)

		switch incrType {
		case "pizzas":
			u.SetPizzas(gs.Resources.Pizzas + incrAmount)

		case "coins":
			u.SetCoins(gs.Resources.Coins + incrAmount)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to perform update: %w", err)
	}

	err = h.sendGameTx(ctx, tx)
	if err != nil {
		return fmt.Errorf("failed to send game tx: %w", err)
	}

	return nil
}
