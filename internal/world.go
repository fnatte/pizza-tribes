package internal

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"github.com/go-redis/redis/v8"
	"google.golang.org/protobuf/encoding/protojson"
)

const WORLD_SIZE = 110
const WORLD_ZONE_SIZE = 10
const WORLD_NUM_ZONES = (WORLD_SIZE / WORLD_ZONE_SIZE) * (WORLD_SIZE / WORLD_ZONE_SIZE)

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
	return y*WORLD_ZONE_SIZE + x
}

func getIdx(x, y int) (zidx, eidx int) {
	zy := y / WORLD_ZONE_SIZE
	zx := x / WORLD_ZONE_SIZE
	zidx = zy*(WORLD_SIZE/WORLD_ZONE_SIZE) + zx
	eidx = (y%WORLD_ZONE_SIZE)*WORLD_ZONE_SIZE + (x % WORLD_ZONE_SIZE)
	return
}

func getZoneIdx(x, y int) int {
	return (y/WORLD_ZONE_SIZE)*(WORLD_SIZE/WORLD_ZONE_SIZE) + (x / WORLD_ZONE_SIZE)
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

func isZoneOpen(z *WorldZone) bool {
	return countTowns(z) < 8
}

type WorldService struct {
	r RedisClient
}

func NewWorldService(r RedisClient) *WorldService {
	return &WorldService{r: r}
}

func (s *WorldService) setEntryXY(ctx context.Context, x, y int, e *WorldEntry) error {
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
		zidxes, err := s.r.ZRange(ctx, "world:open_zones", 0, 0).Result()
		if err != nil {
			return 0, 0, err
		}
		if len(zidxes) != 1 {
			return 0, 0, errors.New("no open zones")
		}
		zidx, err = strconv.Atoi(zidxes[0])
		if err != nil {
			return 0, 0, err
		}

		// Get random entry
		ex := rand.Intn(WORLD_ZONE_SIZE)
		ey := rand.Intn(WORLD_ZONE_SIZE)
		eidx = getEntryIdx(ex, ey)

		// Convert from index to x,y
		zy := zidx / (WORLD_SIZE/WORLD_ZONE_SIZE)
		zx := zidx % (WORLD_SIZE/WORLD_ZONE_SIZE)
		x = zx * (WORLD_ZONE_SIZE) + ex
		y = zy * (WORLD_ZONE_SIZE) + ey

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

	// Close zone if it is fully populated
	zone, err := s.GetZoneIdx(ctx, zidx)
	if err != nil {
		return 0, 0, err
	}
	if !isZoneOpen(zone) {
		s.closeZone(ctx, zidx)
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

func (s *WorldService) GetEntryXY(ctx context.Context, x, y int) (*WorldEntry, error) {
	var zone *WorldZone
	var err error

	zidx, eidx := getIdx(x, y)
	if zone, err = s.GetZoneIdx(ctx, zidx); err != nil {
		return nil, err
	}

	if eidx >= len(zone.Entries) {
		return nil, errors.New("entry not found")
	}

	return zone.Entries[eidx], nil
}

func (s *WorldService) closeZone(ctx context.Context, zidx int) error {
	return s.r.ZRem(ctx, fmt.Sprintf("world:zone:%d", zidx)).Err()
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
	w := WORLD_SIZE / WORLD_ZONE_SIZE
	h := WORLD_SIZE / WORLD_ZONE_SIZE

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			idx := y*WORLD_ZONE_SIZE + x
			b, err := protojson.Marshal(&WorldZone{
				Entries: make([]*WorldEntry, WORLD_ZONE_SIZE*WORLD_ZONE_SIZE),
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

	// Populate open zones. A zone will only be opened if there are no towns in it.
	// Loop through each zone and set its score to its distance from the center
	// This makes the zones closest to the center to be filled first.
	cx := WORLD_SIZE / WORLD_ZONE_SIZE / 2
	cy := WORLD_SIZE / WORLD_ZONE_SIZE / 2
	for x := 0; x < WORLD_SIZE / WORLD_ZONE_SIZE; x++ {
		for y := 0; y < WORLD_SIZE / WORLD_ZONE_SIZE; y++ {
			zidx, _ := getIdx(x * WORLD_ZONE_SIZE, y * WORLD_ZONE_SIZE)
			dx := cx - x
			dy := cy - y
			d := math.Sqrt(float64(dx*dx+dy*dy))
			s.tryOpenZone(ctx, zidx, d)
		}
	}

	return nil
}
