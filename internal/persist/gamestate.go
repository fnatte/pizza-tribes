package persist

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/fnatte/pizza-tribes/internal/models"
	"github.com/fnatte/pizza-tribes/internal/protojson"
	"github.com/fnatte/pizza-tribes/internal/redis"
	"google.golang.org/protobuf/proto"
)

type gameStateRepo struct {
	rdb redis.RedisClient
}

type redisContext struct {
	context.Context
	cmdable redis.Cmdable
}

func NewGameStateRepository(rdb redis.RedisClient) *gameStateRepo {
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
		return nil, fmt.Errorf("failed to get gamestate: %w", err)
	}
	if err == redis.Nil {
		return &models.GameState{}, nil
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
		return nil, fmt.Errorf("failed to unmarshal gamestate: %w", err)
	}

	if len(arr) == 0 {
		return &models.GameState{}, nil
	}

	return arr[0], nil
}

func (r *gameStateRepo) Patch(ctx context.Context, userId string, gs *models.GameState, patchMask *models.PatchMask) error {
	if patchMask == nil {
		return r.Save(ctx, userId, gs)
	}

	gsKey := fmt.Sprintf("user:%s:gamestate", userId)

	pipe := r.rdb.Pipeline()

	for _, p := range patchMask.Paths {
		val, err := models.GetValueByPath(gs, p)
		if err != nil {
			return fmt.Errorf("failed to get value at path %s: %w", p, err)
		}

		jv, err := jsonValue(val)
		if err != nil {
			return fmt.Errorf("failed to get json value: %w", err)
		}

		if val == nil {
			err = redis.RedisJsonDel(pipe, ctx, gsKey, p).Err()
		} else {
			err = redis.RedisJsonSet(pipe, ctx, gsKey, p, jv).Err()
		}
		if err != nil {
			return fmt.Errorf("failed to set %s: %w", p, err)
		}
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("patch failed: %w", err)
	}

	return nil
}

func (r *gameStateRepo) Save(ctx context.Context, userId string, gs *models.GameState) error {
	gsKey := fmt.Sprintf("user:%s:gamestate", userId)

	b, err := protojson.MarshalOptions{
		EmitUnpopulated: true,
	}.Marshal(gs)
	if err != nil {
		return fmt.Errorf("failed to marshal gamestate: %w", err)
	}

	err = r.rdb.JsonSet(ctx, gsKey, ".", string(b)).Err()
	if err != nil {
		return fmt.Errorf("failed to save gamestate: %w", err)
	}

	return nil
}

func jsonValue(v interface{}) (interface{}, error) {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Array || val.Kind() == reflect.Slice {
		mtyp := reflect.TypeOf(new(proto.Message)).Elem()
		if val.Type().Elem().Implements(mtyp) {
			l := val.Len()
			buf := []byte{}
			buf = append(buf, '[')
			for i := 0; i < l; i++ {
				if i != 0 {
					buf = append(buf, ',')
				}

				a := val.Index(i).Interface().(proto.Message)
				b, err := protojson.Marshal(a)
				if err != nil {
					return nil, err
				}
				buf = append(buf, b...)
			}
			buf = append(buf, ']')
			return buf, nil
		}
	}

	if val.Kind() == reflect.Map {
		mtyp := reflect.TypeOf(new(proto.Message)).Elem()
		if val.Type().Elem().Implements(mtyp) {
			m := map[string]*json.RawMessage{}
			iter := val.MapRange()
			for iter.Next() {
				k := iter.Key()
				v := iter.Value().Interface().(proto.Message)
				b, err := protojson.Marshal(v)
				if err != nil {
					return nil, err
				}
				j := json.RawMessage(b)
				m[k.String()] = &j
			}
			return json.Marshal(m)
		}
	}

	switch v := v.(type) {
	case bool:
		if v {
			return "true", nil
		}
		return "false", nil
	case proto.Message:
		return protojson.Marshal(v)
	case string:
		return json.Marshal(v)
	default:
		return v, nil
	}
}
