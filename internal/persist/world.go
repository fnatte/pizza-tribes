package persist

import (
	"context"
	"fmt"

	"github.com/fnatte/pizza-tribes/internal/models"
	"github.com/fnatte/pizza-tribes/internal/protojson"
	"github.com/fnatte/pizza-tribes/internal/redis"
)

type worldRepo struct {
	rc redis.RedisClient
}

func NewWorldRepository(rc redis.RedisClient) *worldRepo {
	return &worldRepo{
		rc: rc,
	}
}

func (r *worldRepo) GetState(ctx context.Context) (*models.WorldState, error) {
	str, err := r.rc.JsonGet(ctx, "world", ".state").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get state: %w", err)
	}

	state := &models.WorldState{}
	protojson.Unmarshal([]byte(str), state)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal world state: %w", err)
	}

	return state, nil
}

