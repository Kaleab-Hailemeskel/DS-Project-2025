package repository

import (
	"context"
	"streaming-service/config"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type IRedisSearchRepo interface {
	IsTokenValid(ctx context.Context, token string) (bool, error)
	MarkTokenAsValid(ctx context.Context, token string, ttl time.Duration) error
}

func InitRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.REDIS_ADDR,
		Password: config.REDIS_PASSWORD,
		DB:       config.REDIS_DB,
	})
	return rdb
}
