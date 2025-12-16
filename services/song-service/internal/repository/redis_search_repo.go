package repository

import (
	"context"
	"fmt"
	"song-service/api/internal/domain"
	"strings"
	"time"

	"github.com/google/uuid"
	redis "github.com/redis/go-redis/v9"
)

const IndexKey = "autocomplete:songs"
const SongKeyPrefix = "song:data:"

type RedisRepository struct {
	// Fields for Redis connection and other configurations
	client *redis.Client
}

func generatePrefixes(s string) []string {
	s = strings.ToLower(s)
	var prefixes []string
	for i := 1; i <= len(s); i++ {
		prefixes = append(prefixes, s[:i])
	}
	return prefixes
}

// DeindexSong implements IRedisSearchRepo.
func (r *RedisRepository) DeindexSong(ctx context.Context, song *domain.Song) error {
	pipe := r.client.Pipeline()

	// 1. Delete the full song data
	songKey := SongKeyPrefix + song.ID.String()
	pipe.Del(ctx, songKey)

	// 2. Generate and remove ALL prefixes from the Sorted Set
	prefixes := generatePrefixes(song.Title)

	// We must remove all the exact members added during indexing
	for _, prefix := range prefixes {
		member := prefix + "::" + song.ID.String()

		// ZREM command to remove the prefix token
		pipe.ZRem(ctx, IndexKey, member)
	}

	// Execute all commands atomically
	_, err := pipe.Exec(ctx)
	return err
}

// IndexSong implements IRedisSearchRepo.
func (r *RedisRepository) IndexSong(ctx context.Context, song *domain.Song) error {
	pipe := r.client.Pipeline()

	// 1. Store the full song data in a Hash or Key-Value store (Essential for retrieval)
	// We use a simple HSET here, storing the ID as a reference value.
	songKey := SongKeyPrefix + song.ID.String()
	pipe.HSet(ctx, songKey, map[string]interface{}{
		"title":     song.Title,
		"artist":    song.Artist,
		"id":        song.ID.String(),
		"album":     song.Album,
		"genre":     song.Genre,
		"image_url": song.ImageURL,
	})

	// 2. Generate and store ALL prefixes in the Sorted Set (The indexing part)
	prefixes := generatePrefixes(song.Title)

	// ZAdd takes Score and Member (the prefix). The score is often used for ranking.
	// We use 0 here for simplicity, focusing only on the prefix match.
	for _, prefix := range prefixes {
		// The member is typically the prefix PLUS the song key,
		// allowing us to look up the full data later.
		member := prefix + "::" + song.ID.String()

		// ZADD command to add the prefix token
		pipe.ZAdd(ctx, IndexKey, redis.Z{
			Score:  float64(time.Now().Unix()), // A score could be the current time for freshness
			Member: member,
		})
	}

	// Execute all commands atomically
	_, err := pipe.Exec(ctx)
	return err
}

// SearchSongsByTitlePrefix implements IRedisSearchRepo.
func (r *RedisRepository) SearchSongsByTitlePrefix(ctx context.Context, titlePrefix string, offset int64, limit int64) ([]*domain.Song, error) {
	if titlePrefix == "" || limit <= 0 {
		return nil, nil
	}

	// 1. --- FAST PREFIX SEARCH (ZRANGEBYLEX) ---

	// Define the titlePrefix range for the sorted set
	start := fmt.Sprintf("[%s", strings.ToLower(titlePrefix))
	end := fmt.Sprintf("[%s\xff", strings.ToLower(titlePrefix))

	// Use ZRANGEBYLEX to find all matching members (tokens + IDs)
	// The Offset and Count parameters enable pagination.
	members, err := r.client.ZRangeByLex(ctx, IndexKey, &redis.ZRangeBy{
		Min:    start,
		Max:    end,
		Offset: offset, // Pagination: Start position
		Count:  limit,  // Pagination: Number of results
	}).Result()

	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("failed to execute ZRANGEBYLEX for search: %w", err)
	}
	if len(members) == 0 {
		return nil, nil
	}

	// Prepare lists for batch fetching
	var songKeys []string

	// --- Parse IDs and Prepare Batch Operations ---
	for _, member := range members {
		parts := strings.Split(member, "::")
		if len(parts) == 2 {
			songIDStr := parts[1]
			songKeys = append(songKeys, SongKeyPrefix+songIDStr)
		}
	}

	// 2. --- BATCH DATA FETCH (HMGET) ---
	pipe := r.client.Pipeline()

	// Fetch the full map data for all keys

	var cmds []*redis.MapStringStringCmd //! changed it from []*redis.StringStringMapCmd to []*redis.MapStringStringCmd, if there is an issue change it back
	for _, key := range songKeys {
		// Queue the retrieval command for each song key
		cmds = append(cmds, pipe.HGetAll(ctx, key))
	}

	// Execute the batch retrieval operations atomically
	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to execute batch retrieval pipeline: %w", err)
	}

	// 3. --- MAP RESULTS TO STRUCTS ---
	songs := make([]*domain.Song, 0, len(cmds))
	for _, cmd := range cmds {
		result, err := cmd.Result()
		if err != nil || len(result) == 0 {
			continue // Skip if data for this ID is missing
		}

		// Map the Hash result back to the domain.Song struct
		id, _ := uuid.Parse(result["id"])

		songs = append(songs, &domain.Song{
			ID:       id,
			Title:    result["title"],
			Artist:   result["artist"],
			Album:    result["album"],
			Genre:    result["genre"],
			ImageURL: result["image_url"],
		})
	}

	return songs, nil
}

func NewRedisRepository(client *redis.Client) IRedisSearchRepo {
	return &RedisRepository{
		client: client,
	}
}
