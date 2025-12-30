package repository

import (
	"context"
	"fmt"
	"song-service/api/config"
	"song-service/api/internal/domain"
	"strings"
	"time"

	"github.com/google/uuid"
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

func (r *RedisRepository) DeindexSong(ctx context.Context, song *domain.Song) error {
	pipe := r.client.Pipeline()
	songKey := SongKeyPrefix + song.ID.String()
	pipe.Del(ctx, songKey)

	prefixes := generatePrefixes(song.Title)
	for _, prefix := range prefixes {
		member := prefix + "::" + song.ID.String()
		pipe.ZRem(ctx, IndexKey, member)
	}

	_, err := pipe.Exec(ctx)
	return err
}

func (r *RedisRepository) IndexSong(ctx context.Context, song *domain.Song) error {
	pipe := r.client.Pipeline()
	songKey := SongKeyPrefix + song.ID.String()

	pipe.HSet(ctx, songKey, map[string]interface{}{
		"title":          song.Title,
		"artist":         song.Artist,
		"id":             song.ID.String(),
		"album":          song.Album,
		"genre":          song.Genre,
		"cover_art_blob": song.CoverArtBlob,
	})

	prefixes := generatePrefixes(song.Title)

	//! comment the following if the search is talking longer than expected
	prefixes = append(prefixes, generatePrefixes(song.Album)...)
	prefixes = append(prefixes, generatePrefixes(song.Genre)...)
	prefixes = append(prefixes, generatePrefixes(song.Artist)...)

	for _, prefix := range prefixes {
		member := prefix + "::" + song.ID.String()
		pipe.ZAdd(ctx, IndexKey, redis.Z{
			Score:  float64(time.Now().Unix()),
			Member: member,
		})
	}

	_, err := pipe.Exec(ctx)
	return err
}

func (r *RedisRepository) SearchSongsByPrefix(ctx context.Context, titlePrefix string, pageNumber int64, pageLimit int64) ([]*domain.Song, error) {
	if titlePrefix == "" || pageLimit <= 0 {
		return nil, nil
	}

	start := fmt.Sprintf("[%s", strings.ToLower(titlePrefix))
	end := fmt.Sprintf("[%s\xff", strings.ToLower(titlePrefix))
	offset := (pageNumber - 1) * min(config.MAX_PAGE_SIZE, pageLimit)

	members, err := r.client.ZRangeByLex(ctx, IndexKey, &redis.ZRangeBy{
		Min:    start,
		Max:    end,
		Offset: offset,
		Count:  pageLimit,
	}).Result()

	if err != nil {
		return nil, fmt.Errorf("failed to execute ZRANGEBYLEX: %w", err)
	}

	if len(members) == 0 {
		return nil, nil
	}

	var songKeys []string
	for _, member := range members {
		parts := strings.Split(member, "::")
		if len(parts) == 2 {
			songKeys = append(songKeys, SongKeyPrefix+parts[1])
		}
	}

	pipe := r.client.Pipeline()
	// Changed from MapStringInterfaceCmd to MapStringStringCmd
	// HGetAll in go-redis returns a map[string]string result.
	var cmds []*redis.MapStringStringCmd
	for _, key := range songKeys {
		cmds = append(cmds, pipe.HGetAll(ctx, key))
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to execute batch retrieval: %w", err)
	}

	songs := make([]*domain.Song, 0, len(cmds))
	uniqueSongIds := make(map[uuid.UUID]struct{})

	for _, cmd := range cmds {
		result, err := cmd.Result()
		if err != nil || len(result) == 0 {
			continue
		}

		id, _ := uuid.Parse(result["id"])

		if _, exists := uniqueSongIds[id]; exists {
			continue
		}
		uniqueSongIds[id] = struct{}{}

		// No interface assertion needed now as result is map[string]string
		songs = append(songs, &domain.Song{
			ID:           id,
			Title:        result["title"],
			Artist:       result["artist"],
			Album:        result["album"],
			Genre:        result["genre"],
			CoverArtBlob: []byte(result["cover_art_blob"]),
		})
	}

	return songs, nil
}

func NewRedisRepository(client *redis.Client) IRedisSearchRepo {
	return &RedisRepository{
		client: client,
	}
}
