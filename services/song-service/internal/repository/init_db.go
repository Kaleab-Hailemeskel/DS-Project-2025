package repository

import (
	"context"
	"fmt"
	"song-service/api/config"
	"song-service/api/internal/domain"

	"github.com/google/uuid"
	redis "github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ISongRepo interface {
	GetSong(id uuid.UUID) (*domain.Song, error)
	SaveSong(song *domain.Song) (*domain.Song, error)
	GetAllSongs(mulicListPerPage, pageNumber int) ([]*domain.Song, error)
	GetSongByArtist(artist string) ([]*domain.Song, error)
	GetSongByTitle(title string) ([]*domain.Song, error)
	GetSongByAlbum(album string) ([]*domain.Song, error)
	GetSongByGenre(genre string) ([]*domain.Song, error)
	UpdateSong(song *domain.Song) (*domain.Song, error)
	DeleteSong(id uuid.UUID) error
}

type IRedisSearchRepo interface {
	IndexSong(ctx context.Context, song *domain.Song) error
	DeindexSong(ctx context.Context, song *domain.Song) error
	SearchSongsByTitlePrefix(ctx context.Context, titlePrefix string, offset, limit int64) ([]*domain.Song, error)
	//* they aren't implemented yet
	/*
		SearchSongsByArtist(ctx context.Context, artist string) ([]*domain.Song, error)
		SearchSongsByAlbum(ctx context.Context, album string) ([]*domain.Song, error)
		SearchSongsByGenre(ctx context.Context, genre string) ([]*domain.Song, error)
	*/
}

func InitRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.REDIS_ADDR,
		Password: config.REDIS_PASSWORD, // no password set
		DB:       config.REDIS_DB,       // use default DB
	})
	return rdb
}
func InitPostgresDB() *gorm.DB {
	// Implementation for initializing PostgreSQL DB connection

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		config.POSTGRES_HOST,
		config.POSTGRES_USER,
		config.POSTGRES_PASSWORD,
		config.POSTGRES_DB,
		config.POSTGRES_PORT,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	return db
}
