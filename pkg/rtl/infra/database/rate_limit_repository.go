package database

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimitRedisRepository struct {
	RDB *redis.Client
}

func NewRateLimitRedisRepository(rdb *redis.Client) *RateLimitRedisRepository {
	return &RateLimitRedisRepository{
		RDB: rdb,
	}
}

func (r *RateLimitRedisRepository) GetByKey(ctx context.Context, key string) (int, error) {
	counter, err := r.RDB.Get(ctx, key).Int()
	if err != nil {
		return 0, err
	}
	return counter, nil
}

func (r *RateLimitRedisRepository) Add(ctx context.Context, key string) error {
	_, err := r.RDB.Incr(ctx, key).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *RateLimitRedisRepository) SetTimeOutByKey(ctx context.Context, key string, timeSec int) error {
	err := r.RDB.Set(ctx, key, 1, time.Duration(timeSec)*time.Second).Err()
	if err != nil {
		return err
	}
	return nil
}
