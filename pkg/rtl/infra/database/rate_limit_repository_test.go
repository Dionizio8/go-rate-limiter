package database

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

const redis_addr = "localhost:6379"

var rdb *redis.Client

func TestMain(m *testing.M) {
	rdb = redis.NewClient(&redis.Options{
		Addr: redis_addr,
		DB:   0,
	})
}

func TestRateLimitRepository_GetByKey_Ok(t *testing.T) {
	err := rdb.Set(context.Background(), "teste", 1, time.Duration(2)*time.Second).Err()
	if err != nil {
		t.Error(err)
	}

	rtlRedisRepository := NewRateLimitRedisRepository(rdb)
	counter, err := rtlRedisRepository.GetByKey(context.Background(), "teste")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 1, counter)
}

func TestRateLimitRepository_GetByKey_NotFount(t *testing.T) {
	rtlRedisRepository := NewRateLimitRedisRepository(rdb)
	counter, err := rtlRedisRepository.GetByKey(context.Background(), "teste_not_found")
	if err != nil && err != redis.Nil {
		t.Error(err)
	}
	assert.Equal(t, 0, counter)
	assert.Equal(t, redis.Nil, err)
}

func TestRateLimitRepository_Add_Ok(t *testing.T) {
	rtlRedisRepository := NewRateLimitRedisRepository(rdb)
	rtlRedisRepository.Add(context.Background(), "teste_add_ok")

	counter, err := rtlRedisRepository.GetByKey(context.Background(), "teste_add_ok")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 1, counter)
}

func TestRateLimitRepository_SetTimeOutByKey_Ok(t *testing.T) {
	rtlRedisRepository := NewRateLimitRedisRepository(rdb)
	rtlRedisRepository.SetTimeOutByKey(context.Background(), "teste_set_timeout_ok", 2)

	counter, err := rtlRedisRepository.GetByKey(context.Background(), "teste_add_ok")
	assert.Nil(t, err)
	assert.Equal(t, 1, counter)

	time.Sleep(2 * time.Second)

	counter, err = rtlRedisRepository.GetByKey(context.Background(), "teste_add_ok")
	assert.Nil(t, err)
	assert.Equal(t, 0, counter)
}
