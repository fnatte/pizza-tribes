package main

import (
	"context"

	"github.com/fnatte/pizza-tribes/internal/game"
	"github.com/fnatte/pizza-tribes/internal/game/models"
)

func (h *handler) handleReadReport(ctx context.Context, userId string, m *models.ClientMessage_ReadReport) error {
	return game.MarkReportAsRead(ctx, h.rdb, userId, m.Id)
}
