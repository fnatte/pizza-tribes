package internal

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/encoding/protojson"
)

const WORLD_WIDTH = 110
const WORLD_HEIGHT = 110
const WORLD_ZONE_WIDTH = 10
const WORLD_ZONE_HEIGHT = 10

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

func countTowns(z *WorldZone) int {
	count := 0
	for i := range z.Entries {
		if z.Entries[i].GetTown() != nil {
			count++
		}
	}
	return count
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

func (s *WorldService) tryOpenZone(ctx context.Context, zidx int, score float64) error {
	// Check middle zone
	zone, err := s.GetZoneIdx(ctx, zidx)
	if err != nil {
		return err
	}

	if zone == nil {
		return nil
	}

	count := countTowns(zone)
	if count == 0 {
		s.r.ZAdd(ctx, "world:open_zones", &redis.Z{
			Score:  score,
			Member: zidx,
		}).Err()
	}

	return nil
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

	// Populate open zones
	allOpenZones, err := s.r.ZRange(ctx, "world:open_zones", -1, -1).Result()
	if err != nil {
		return err
	}
	if len(allOpenZones) == 0 {
		// Walk the edges from the center and outwards
		// Visualized JS-version: https://codepen.io/fnatteh/pen/MWpYrKP
		cx := WORLD_WIDTH / 2
		cy := WORLD_HEIGHT / 2
		dx := 10
		dy := 0
		x := cx - 10
		y := cy - 10

		zidx, _ := getIdx(cx, cy)
		s.tryOpenZone(ctx, zidx, 0)

		for r := 0; r < 5; r++ {
			for side := 0; side < 4; side++ {
				for i := 0; i < (r+1)*2; i++ {
					zidx, _ = getIdx(x, y)
					s.tryOpenZone(ctx, zidx, float64(r*len(xyOffsets)*i))

					x = x + dx
					y = y + dy
				}

				// turn
				t := dx
				dx = -dy
				dy = t
			}

			y = y - 10
			x = x - 10
			dx = 10
			dy = 0
		}
	}

	return nil
}
