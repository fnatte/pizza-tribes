package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

type Reminder struct {
	Id       string        `json:"id"`
	Interval time.Duration `json:"interval"`
	Offset   time.Duration `json:"offset"`
}

func (r *Reminder) NextOccurenceAfter(t time.Time) time.Time {
	next := t.Truncate(r.Interval).Add(r.Offset)
	if next.Before(t) {
		next = next.Add(r.Interval)
	}
	return next
}

func (r *Reminder) NextOccurence() time.Time {
	return r.NextOccurenceAfter(time.Now())
}

func ScheduleReminder(ctx context.Context, rc RedisClient, r *Reminder) (int64, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return 0, err
	}

	err = rc.Set(ctx, fmt.Sprintf("reminder:%s", r.Id), b, 0).Err()
	if err != nil {
		return 0, err
	}
	t := r.NextOccurence().UnixNano()

	err = rc.ZAdd(ctx, "reminders", &redis.Z{
		Score:  float64(t),
		Member: r.Id,
	}).Err()
	if err != nil {
		return 0, err
	}

	return t, nil
}

func NextReminder(ctx context.Context, rc RedisClient) (*Reminder, error) {
	packed, err := rc.ZRangeWithScores(ctx, "reminders", 0, 0).Result()
	if err != nil {
		return nil, err
	}

	if len(packed) == 0 {
		return nil, nil
	}

	timestamp := int64(packed[0].Score)
	if timestamp > time.Now().UnixNano() {
		return nil, nil
	}

	member, ok := packed[0].Member.(string)
	if !ok {
		return nil, errors.New("failed to read reminders member")
	}

	jsonStr, err := rc.Get(ctx, fmt.Sprintf("reminder:%s", member)).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get reminder: %w", err)
	}

	r := &Reminder{}
	err = json.Unmarshal([]byte(jsonStr), r)
	if err != nil {
		return nil, fmt.Errorf("failed to parse reminder: %w", err)
	}

	removed, err := rc.ZRem(ctx, "reminders", member).Result()
	if err != nil {
		return nil, err
	}

	if removed != 1 {
		return nil, nil
	}

	_, err = ScheduleReminder(ctx, rc, r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func HandleReminders(ctx context.Context, rc RedisClient, h func(r *Reminder)) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			r, err := NextReminder(ctx, rc)
			if err != nil {
				log.Error().Err(err).Msg("Failed to get next reminder")
				time.Sleep(1 * time.Second)
				continue
			}

			if r == nil {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			h(r)
		}
	}
}
