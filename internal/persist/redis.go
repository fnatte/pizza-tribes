package persist

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fnatte/pizza-tribes/internal"
	"github.com/fnatte/pizza-tribes/internal/models"
	"github.com/fnatte/pizza-tribes/internal/protojson"
	"github.com/go-redis/redis/v8"
)

type gameStateRepo struct {
	rdb internal.RedisClient
}

func NewGameStateRepository(rdb internal.RedisClient) *gameStateRepo {
	return &gameStateRepo{
		rdb: rdb,
	}
}

func (r *gameStateRepo) NewMutex(userId string) Mutex {
	gsKey := fmt.Sprintf("user:%s:gamestate", userId)

	return r.rdb.NewMutex("lock:" + gsKey)
}

func (r *gameStateRepo) Get(ctx context.Context, userId string) (*models.GameState, error) {
	gsKey := fmt.Sprintf("user:%s:gamestate", userId)
	s, err := r.rdb.JsonGet(ctx, gsKey, "$").Result()

	if err != nil && err != redis.Nil {
		return nil, err
	}

	var arr []*models.GameState

	if err = protojson.UnmarshalArray([]byte(s), func(buf json.RawMessage) error {
		gs := &models.GameState{}
		err := protojson.Unmarshal(buf, gs)
		if err != nil {
			return err
		}

		arr = append(arr, gs)

		return nil
	}); err != nil {
		return nil, err
	}

	if len(arr) == 0 {
		return &models.GameState{}, nil
	}

	return arr[0], nil
}

func (r *gameStateRepo) Patch(ctx context.Context, userId string, patch *Patch) error {
	gsKey := fmt.Sprintf("user:%s:gamestate", userId)

	pipe := r.rdb.Pipeline()

	for _, op := range patch.Ops {
		path := JsonPointerToJsonPath(op.Path)
		switch op.Op {
		case "replace":
			err := internal.RedisJsonSet(pipe, ctx, gsKey, path, jsonValue(op.Value)).Err()
			if err != nil {
				return fmt.Errorf("failed to set %s: %w", op.Path, err)
			}
		}
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("patch failed: %w", err)
	}

	return nil
}

func jsonValue(v interface{}) interface{} {
	switch v := v.(type) {
	case bool:
		if v {
			return "true"
		}
		return "false"
	default:
		return v
	}
}

