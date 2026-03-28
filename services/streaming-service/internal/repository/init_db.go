package repository

import (
	"context"
	"log"
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
	log.Printf("Initializing Redis client with config: Addr=%s, Password=%s, DB=%d",
		config.REDIS_ADDR, config.REDIS_PASSWORD, config.REDIS_DB)

	// Context is required for go-redis v8+
	ctx := context.Background()

	// Perform the Ping
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("💀❌ Could not connect to Redis: %v", err)
	}

	return rdb
}
