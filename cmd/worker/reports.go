package main

import (
	"context"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/fnatte/pizza-tribes/internal/models"
)

func (h *handler) handleReadReport(ctx context.Context, userId string, m *models.ClientMessage_ReadReport) error {
	return internal.MarkReportAsRead(ctx, h.rdb, userId, m.Id)
}
