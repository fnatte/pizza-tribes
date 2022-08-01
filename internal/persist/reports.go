package persist

import (
	"context"
	"fmt"

	"github.com/fnatte/pizza-tribes/internal/models"
	"github.com/fnatte/pizza-tribes/internal/protojson"
	"github.com/fnatte/pizza-tribes/internal/redis"
)

type reportsRepo struct {
	rdb redis.RedisClient
}

func NewReportsRepository(rdb redis.RedisClient) *reportsRepo {
	return &reportsRepo{
		rdb: rdb,
	}
}

func (r *reportsRepo) Save(ctx context.Context, userId string, report *models.Report) error {
	b, err := protojson.Marshal(report)
	if err != nil {
		return fmt.Errorf("failed to marshal report: %w", err)
	}

	reportsKey := fmt.Sprintf("user:%s:reports", userId)
	reportIndexKey := fmt.Sprintf("user:%s:reportsByDate", userId)

	pipe := r.rdb.TxPipeline()
	pipe.HSet(ctx, reportsKey, report.Id, b)
	pipe.ZAdd(ctx, reportIndexKey, &redis.Z{
		Score:  float64(report.CreatedAt),
		Member: report.Id,
	})

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to save report: %w", err)
	}

	return nil
}

func (r *reportsRepo) Get(ctx context.Context, userId string) ([]*models.Report, error) {
	reportsKey := fmt.Sprintf("user:%s:reports", userId)
	reportIndexKey := fmt.Sprintf("user:%s:reportsByDate", userId)

	ids, err := r.rdb.ZRevRange(ctx, reportIndexKey, 0, 10).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to read from reports index: %s", err)
	}
	if len(ids) == 0 {
		return []*models.Report{}, nil
	}
	res, err := r.rdb.HMGet(ctx, reportsKey, ids...).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to read reports: %s", err)
	}

	reports := make([]*models.Report, len(res))
	for i, b := range res {
		reports[i] = &models.Report{}
		err = protojson.Unmarshal([]byte(b.(string)), reports[i])
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal report: %s", err)
		}
	}

	return reports, nil
}

func (r *reportsRepo) MarkRead(ctx context.Context, userId string, reportId string) error {
	report := &models.Report{}

	// Get report
	reportsKey := fmt.Sprintf("user:%s:reports", userId)
	str, err := r.rdb.HGet(ctx, reportsKey, reportId).Result()
	if err != nil {
		return fmt.Errorf("failed to get report: %s", err)
	}
	err = protojson.Unmarshal([]byte(str), report)

	// Mark report as read
	report.Unread = false

	// Write back report
	b, err := protojson.Marshal(report)
	if err != nil {
		return fmt.Errorf("failed to marshal report: %w", err)
	}
	return r.rdb.HSet(ctx, reportsKey, reportId, b).Err()
}
