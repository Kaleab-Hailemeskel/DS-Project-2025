package repository

import (
	"context"
	"strings"
	"time"

	redis "github.com/redis/go-redis/v9"
)

const IndexKey = "autocomplete:songs"
const SongKeyPrefix = "song:data:"

type RedisRepository struct {
	client *redis.Client
}

// IsTokenValid implements [IRedisSearchRepo].
func (r *RedisRepository) IsTokenValid(ctx context.Context, token string) (bool, error) {
	// Exists returns 1 if key exists, 0 if not
	val, err := r.client.Exists(ctx, "auth_token:"+token).Result()
	if err != nil {
		return false, err
	}
	return val > 0, nil
}

// MarkTokenAsValid implements [IRedisSearchRepo].
func (r *RedisRepository) MarkTokenAsValid(ctx context.Context, token string, ttl time.Duration) error {
	// We store a simple "1" as the value since we only care about key existence
	return r.client.Set(ctx, "auth_token:"+token, "1", ttl).Err()
}

func generatePrefixes(s string) []string {
	s = strings.ToLower(s)
	var prefixes []string
	for i := 1; i <= len(s); i++ {
		prefixes = append(prefixes, s[:i])
	}
	return prefixes
}

func NewRedisRepository(client *redis.Client) IRedisSearchRepo {
	return &RedisRepository{
		client: client,
	}
}
