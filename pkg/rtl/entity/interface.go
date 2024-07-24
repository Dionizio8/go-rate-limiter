package entity

import "context"

type RateLimitRepositoryInterface interface {
	GetByKey(ctx context.Context, key string) (int, error)
	Add(ctx context.Context, key string) error
	SetTimeOutByKey(ctx context.Context, key string, timeSec int) error
}
