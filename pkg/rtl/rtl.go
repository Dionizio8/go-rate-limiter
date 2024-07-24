package rtl

import (
	"context"
	"fmt"

	"github.com/Dionizio8/go-rate-limiter/pkg/rtl/entity"
	"github.com/redis/go-redis/v9"
)

const (
	prefixKey = "GORTL"
)

type RTL struct {
	RTLRepository entity.RateLimitRepositoryInterface
	IPLimit       int
	BlockTime     int
	ListToken     map[string]int // [token]limit
}

func NewRTL(rtlRepository entity.RateLimitRepositoryInterface, ipLimit int, blockTime int) *RTL {
	return &RTL{
		RTLRepository: rtlRepository,
		IPLimit:       ipLimit,
		BlockTime:     blockTime,
		ListToken:     make(map[string]int),
	}
}

func (r *RTL) SetToken(token string, limitTime int) {
	r.ListToken[token] = limitTime
}

func (r *RTL) Validate(ctx context.Context, key string) (bool, error) {
	dbKey := r.getKey(key)
	counter, err := r.RTLRepository.GetByKey(ctx, dbKey)
	if err == redis.Nil {
		timeSec := r.getTime(key)
		err = r.RTLRepository.SetTimeOutByKey(ctx, dbKey, timeSec)
		if err != nil {
			return false, err
		}
		return true, nil
	} else if err != nil {
		return false, err
	}

	if counter >= r.IPLimit {
		return false, nil
	}

	err = r.RTLRepository.Add(ctx, dbKey)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *RTL) getTime(token string) int {
	var timeSec int = r.BlockTime
	for t, l := range r.ListToken {
		if t == token {
			timeSec = l
			break
		}
	}
	return timeSec
}

func (r *RTL) getKey(token string) string {
	return fmt.Sprintf("%s_%s", prefixKey, token)
}
