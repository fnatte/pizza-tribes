package internal

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fnatte/pizza-tribes/internal/models"
	. "github.com/fnatte/pizza-tribes/internal/models"
	"github.com/fnatte/pizza-tribes/internal/protojson"
	"github.com/fnatte/pizza-tribes/internal/spot_finder"
	"github.com/go-redis/redis/v8"
)

const WORLD_SIZE = 110

type xy struct{ x, y int }

var xyOffsets []xy = []xy{
	{x: 1, y: -1},
	{x: 1, y: 0},
	{x: 1, y: 1},
	{x: 0, y: 1},
	{x: -1, y: 1},
	{x: -1, y: 0},
	{x: -1, y: -1},
	{x: 0, y: -1},
}

type WorldService struct {
	r RedisClient
}

func parseWorldKey(key string) (int, int, error) {
	split := strings.Split(key, ":")
	if len(split) != 2 {
		return 0, 0, fmt.Errorf("unexpected key format: %s", key)
	}

	x, err := strconv.ParseInt(split[0], 10, 32)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse x as integer: %w", err)
	}

	y, err := strconv.ParseInt(split[1], 10, 32)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse y as integer: %w", err)
	}

	return int(x), int(y), nil
}

func getWorldKey(x, y int) string {
	return fmt.Sprintf("%d:%d", x, y)
}

func NewWorldService(r RedisClient) *WorldService {
	return &WorldService{r: r}
}

func (s *WorldService) Start(ctx context.Context) error {
	state, err := s.GetState(ctx)
	if err != nil {
		return fmt.Errorf("failed to start world: %w", err)
	}

	state.Type = &models.WorldState_Started_{}

	b, err := protojson.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to start state: %w", err)
	}

	return s.r.JsonSet(ctx, "world", ".state", b).Err()
}

func (s *WorldService) End(ctx context.Context, winnerUserId string) (*WorldState, error) {
	state, err := s.GetState(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to end world: %w", err)
	}

	state.Type = &models.WorldState_Ended_{
		Ended: &WorldState_Ended{
			WinnerUserId: winnerUserId,
		},
	}

	b, err := protojson.Marshal(state)
	if err != nil {
		return nil, fmt.Errorf("failed to end state: %w", err)
	}

	if err := s.r.JsonSet(ctx, "world", ".state", b).Err(); err != nil {
		return nil, err
	}

	return state, nil
}

func (s *WorldService) GetState(ctx context.Context) (*WorldState, error) {
	str, err := s.r.JsonGet(ctx, "world", ".state").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get state: %w", err)
	}

	state := &WorldState{}
	protojson.Unmarshal([]byte(str), state)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal world state: %w", err)
	}

	return state, nil
}

func (s *WorldService) setEntryXY(ctx context.Context, x, y int, e *WorldEntry) error {
	data, err := protojson.Marshal(e)
	if err != nil {
		return err
	}

	path := fmt.Sprintf(".entries[\"%d:%d\"]", x, y)

	return s.r.JsonSet(ctx, "world", path, data).Err()
}

func (s *WorldService) GetEntryXY(ctx context.Context, x, y int) (*WorldEntry, error) {
	path := fmt.Sprintf(".entries[\"%d:%d\"]", x, y)

	str, err := s.r.JsonGet(ctx, "world", path).Result()
	if err != nil {
		if err == redis.Nil || IsRedisJsonPathDoesNotExistError(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get entry at %d:%d: %w", x, y, err)
	}

	entry := &WorldEntry{}
	protojson.Unmarshal([]byte(str), entry)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal world entry: %w", err)
	}

	return entry, nil
}

func (s *WorldService) GetEntries(ctx context.Context, x, y, radius int) (map[string]*WorldEntry, error) {
	str, err := s.r.JsonGet(ctx, "world", ".").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get entries: %w", err)
	}

	world := &World{}
	protojson.Unmarshal([]byte(str), world)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal world entry: %w", err)
	}

	entries := map[string]*WorldEntry{}

	r2 := radius * radius

	var ex, ey int

	for key := range world.Entries {
		if ex, ey, err = parseWorldKey(key); err != nil {
			return nil, fmt.Errorf("failed to unmarshal world entry key: %w", err)
		}

		if (ex-x)*(ex-x)+(ey-y)*(ey-y) < r2 {
			entries[key] = world.Entries[key]
		}
	}

	return entries, nil
}

func (s *WorldService) AcquireTown(ctx context.Context, userId string) (x, y int, err error) {
	// Find a spot for the new town
	entries, err := s.GetEntries(ctx, WORLD_SIZE/2, WORLD_SIZE/2, WORLD_SIZE)
	if err != nil {
		return 0, 0, err
	}
	v2s := make([]spot_finder.Vec2, 0, len(entries))
	for k, _ := range entries {
		ex, ey, err := parseWorldKey(k)
		if err != nil {
			return 0, 0, err
		}
		v2s = append(v2s, spot_finder.NewVec2(float64(ex), float64(ey)))
	}
	x, y = spot_finder.FindSpotForNewTown(v2s)

	// Verify that the spot is empty
	if entries[getWorldKey(x, y)] != nil {
		return 0, 0, fmt.Errorf("spot_finder returned a spot that was not empty")
	}

	// Assign a town to the world entry
	err = s.setEntryXY(ctx, x, y, &WorldEntry{
		Object: &WorldEntry_Town_{
			Town: &WorldEntry_Town{
				UserId: userId,
			},
		},
	})
	if err != nil {
		return 0, 0, err
	}

	// Update the user game state with coordinates to its town.
	// Since the user will often need the coords of its of town,
	// we don't want to search the entire world space.
	err = s.r.JsonSet(ctx, fmt.Sprintf("user:%s:gamestate", userId), ".townX", x).Err()
	if err != nil {
		return 0, 0, err
	}
	err = s.r.JsonSet(ctx, fmt.Sprintf("user:%s:gamestate", userId), ".townY", y).Err()
	if err != nil {
		return 0, 0, err
	}

	return
}

func (s *WorldService) Initialize(ctx context.Context) error {
	if s.r.Exists(ctx, "world").Val() == 0 {
		world := models.World{
			Entries: map[string]*models.WorldEntry{},
			State: &models.WorldState{
				Type:      &WorldState_Starting_{},
				StartTime: time.Now().Truncate(24 * time.Hour).Add(36 * time.Hour).Unix(),
			},
		}

		b, err := protojson.MarshalWithUnpopulated(&world)
		if err != nil {
			return fmt.Errorf("failed to marshal world: %w", err)
		}
		err = s.r.JsonSet(ctx, "world", ".", b).Err()

		if err != nil {
			return fmt.Errorf("failed to save world: %w", err)
		}
	}

	return nil
}
