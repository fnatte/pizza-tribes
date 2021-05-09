package main

import (
	"context"

	"github.com/fnatte/pizza-tribes/internal"
)

func (h *handler) handleReadReport(ctx context.Context, userId string, m *internal.ClientMessage_ReadReport) error {
	return internal.MarkReportAsRead(ctx, h.rdb, userId, m.Id)
}
