package internal

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/encoding/protojson"
)

const WORLD_WIDTH = 100
const WORLD_HEIGHT = 100
const WORLD_ZONE_WIDTH = 10
const WORLD_ZONE_HEIGHT = 10

func getEntryIdx(x, y int) int {
	return y*WORLD_ZONE_WIDTH + x
}

func getIdx(x, y int) (zidx, eidx int) {
	zy := y / WORLD_ZONE_HEIGHT
	zx := x / WORLD_ZONE_WIDTH
	zidx = zy*(WORLD_WIDTH/WORLD_ZONE_WIDTH) + zx
	eidx = (y%WORLD_ZONE_HEIGHT)*WORLD_ZONE_WIDTH + (x % WORLD_ZONE_WIDTH)
	return
}

func getZoneIdx(x, y int) int {
	return (y/WORLD_ZONE_HEIGHT)*(WORLD_WIDTH/WORLD_ZONE_WIDTH) + (x / WORLD_ZONE_WIDTH)
}

func getZoneKey(idx int) string {
	return fmt.Sprintf("world:zone:%d", idx)
}

type WorldService struct {
	r RedisClient
}

func NewWorldService(r RedisClient) *WorldService {
	return &WorldService{r: r}
}

func (s *WorldService) SetEntryXY(ctx context.Context, x, y int, e *WorldEntry) error {
	data, err := protojson.Marshal(e)
	if err != nil {
		return err
	}

	zidx, eidx := getIdx(x, y)
	path := fmt.Sprintf(".entries[%d]", eidx)

	return s.r.JsonSet(ctx, getZoneKey(zidx), path, data).Err()
}

func (s *WorldService) AcquireTown(ctx context.Context, userId string) (x, y int, err error) {
	var zidx, eidx int

	// Loop until we find a spot
	for {
		x = rand.Intn(WORLD_WIDTH)
		y = rand.Intn(WORLD_HEIGHT)
		zidx, eidx = getIdx(x, y)

		log.Info().Int("zidx", zidx).Msg("GetZoneIdx")
		zone, err := s.GetZoneIdx(ctx, zidx)
		if err != nil {
			return 0, 0, err
		}

		if zone != nil && eidx < len(zone.Entries) {
			e := zone.Entries[eidx]
			if e.GetTown() != nil {
				continue
			}
		}

		break
	}

	err = s.SetEntryXY(ctx, x, y, &WorldEntry{
		Object: &WorldEntry_Town_{
			Town: &WorldEntry_Town{
				UserId: userId,
			},
		},
	})
	if err != nil {
		return 0, 0, err
	}
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

func (s *WorldService) GetZoneXY(ctx context.Context, x, y int) (*WorldZone, error) {
	idx := getZoneIdx(x, y)
	return s.GetZoneIdx(ctx, idx)
}

func (s *WorldService) GetZoneIdx(ctx context.Context, idx int) (*WorldZone, error) {
	b, err := s.r.JsonGet(ctx, getZoneKey(idx), ".").Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	z := &WorldZone{}
	err = protojson.Unmarshal([]byte(b), z)
	if err != nil {
		return nil, err
	}

	return z, nil
}

func (s *WorldService) Initilize(ctx context.Context) error {
	w := WORLD_WIDTH / WORLD_ZONE_WIDTH
	h := WORLD_HEIGHT / WORLD_ZONE_HEIGHT

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			idx := y*WORLD_ZONE_WIDTH + x
			b, err := protojson.Marshal(&WorldZone{
				Entries: make([]*WorldEntry, WORLD_ZONE_HEIGHT*WORLD_ZONE_WIDTH),
			})
			if err != nil {
				return err
			}

			key := getZoneKey(idx)

			if s.r.JsonGet(ctx, key, ".").Err() == redis.Nil {
				s.r.JsonSet(ctx, key, ".", b)
			}
		}
	}

	return nil
}
