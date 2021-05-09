package internal

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"google.golang.org/protobuf/encoding/protojson"
)

type ReportsService struct {
	r RedisClient
}

func NewReportsService(r RedisClient) *ReportsService {
	return &ReportsService{r: r}
}

func SaveReport(ctx context.Context, r redis.Cmdable, userId string, report *Report) error {
	b, err := protojson.Marshal(report)
	if err != nil {
		return fmt.Errorf("failed to marshal report: %w", err)
	}

	reportsKey := fmt.Sprintf("user:%s:reports", userId)
	reportIndexKey := fmt.Sprintf("user:%s:reportsByDate", userId)

	pipe := r.TxPipeline()
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

func GetReports(ctx context.Context, r redis.Cmdable, userId string) ([]*Report, error) {
	reportsKey := fmt.Sprintf("user:%s:reports", userId)
	reportIndexKey := fmt.Sprintf("user:%s:reportsByDate", userId)

	ids, err := r.ZRevRange(ctx, reportIndexKey, 0, 10).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to read from reports index: %s", err)
	}
	res, err := r.HMGet(ctx, reportsKey, ids...).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to read reports: %s", err)
	}

	reports := make([]*Report, len(res))
	for i, b := range res {
		reports[i] = &Report{}
		err = protojson.Unmarshal([]byte(b.(string)), reports[i])
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal report: %s", err)
		}
	}

	return reports, nil
}

func MarkReportAsRead(ctx context.Context, r redis.Cmdable, userId string, reportId string) error {
	report := &Report{}

	// Get report
	reportsKey := fmt.Sprintf("user:%s:reports", userId)
	str, err := r.HGet(ctx, reportsKey, reportId).Result()
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
	return r.HSet(ctx, reportsKey, reportId, b).Err()
}

