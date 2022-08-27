package ws

import (
	"context"

	"github.com/fnatte/pizza-tribes/internal/game"
	"github.com/fnatte/pizza-tribes/internal/game/models"
	"github.com/fnatte/pizza-tribes/internal/game/protojson"
	"github.com/go-redis/redis/v8"
)

func Send(ctx context.Context, r redis.Cmdable, userId string, msg *models.ServerMessage) error {
	b, err := protojson.Marshal(msg)
	if err != nil {
		return err
	}

	return r.RPush(ctx, "wsout", &game.OutgoingMessage{
		ReceiverId: userId,
		Body:       string(b),
	}).Err()
}
