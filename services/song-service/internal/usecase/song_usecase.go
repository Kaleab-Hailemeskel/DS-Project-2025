package usecase

import (
	"context"
	"fmt"
	"log"
	"song-service/api/config"
	"song-service/api/internal/domain"
	"song-service/api/internal/repository"
	"strconv"

	"github.com/google/uuid"
)

type SongUsecase struct {
	songRepo        repository.ISongRepo
	redisSearchRepo repository.IRedisSearchRepo
}

// SaveBlobImage implements [ISongUsecase].
func (s *SongUsecase) SaveBlobImage(id uuid.UUID, blob []byte) error {
	if blob == nil {
		return fmt.Errorf("blob is nil")
	} else if id == uuid.Nil {
		return fmt.Errorf("id is nil")
	}
	return s.songRepo.SaveBlobImage(id, blob)
}

// GetAllSong implements [ISongUsecase].
func (s *SongUsecase) GetAllSong(pageNumber string, pageLimit string) ([]*domain.Song, error) {
	pageLimitInt := int(config.MAX_PAGE_SIZE)
	pageNumberInt := 0
	if value, err := strconv.Atoi(pageLimit); err == nil {
		pageLimitInt = value
	}
	if value, err := strconv.Atoi(pageNumber); err == nil {
		pageNumberInt = value
	}
	return s.songRepo.GetAllSongs(pageLimitInt, pageNumberInt)
}

// DeindexSong implements ISongUsecase.
func (s *SongUsecase) DeindexSong(song *domain.Song) error {
	return s.redisSearchRepo.DeindexSong(context.Background(), song)
}

// FilterSongs implements ISongUsecase.
func (s *SongUsecase) FilterSongs(query string, filters map[string]interface{}) ([]*domain.Song, error) {
	panic("unimplemented")
}

// IndexSong implements ISongUsecase.
func (s *SongUsecase) IndexSong(song *domain.Song) error {
	return s.redisSearchRepo.IndexSong(context.Background(), song)
}

// SearchSongsByPrefix implements ISongUsecase.
func (s *SongUsecase) SearchSongsByPrefix(stringPrefix, pageNumber, pageLimit string) ([]*domain.Song, error) {
	// Convert offset and pageLimit to integers
	pageNumberInt, err := strconv.ParseInt(pageNumber, 10, 64)
	if err != nil {
		return nil, err
	}
	pageLimitInt, err := strconv.ParseInt(pageLimit, 10, 64)
	if err != nil {
		return nil, err
	}

	searchCache, err := s.redisSearchRepo.SearchSongsByPrefix(context.Background(), stringPrefix, pageNumberInt, pageLimitInt)
	if err != nil {
		return nil, err
	}
	log.Println("✅ Result From Redis: ")
	log.Println("------------------------------")
	log.Printf("Songs found from redis cache: %#v\n", searchCache)
	log.Printf("Error: %v\n", err)
	log.Println("------------------------------")
	if searchCache == nil {
		// If no results in cache, search in the song repository (database)

		res, err := s.songRepo.SearchSongs(stringPrefix, int(pageLimitInt), int(pageNumberInt))
		log.Println("✅ Result From Database: ")
		log.Println("------------------------------")
		log.Printf("Songs From Postgres: %#v\n", res)
		log.Printf("Error: %v\n", err)
		log.Println("------------------------------")
		if res != nil {
			log.Printf("song found from database %#v", res)
		} else {
			log.Printf("no song found from database")
		}
		if err != nil {
			return nil, err
		}
		for _, song := range res {
			// Index each song found in the database to Redis for future searches
			_ = s.redisSearchRepo.IndexSong(context.Background(), song)
		}
		return res, nil
	}
	return searchCache, nil
}

// UploadFileToArchive implements ISongUsecase.
func (u *SongUsecase) SaveSongMetaData(songMetaData *domain.Song) (*domain.Song, error) {
	if res, err := u.songRepo.GetOneSongExact(songMetaData.Title, songMetaData.Album, songMetaData.Artist, songMetaData.Genre); res != nil || err != nil {
		return nil, fmt.Errorf("while checking for duplicates for song")
	}
	return u.songRepo.SaveSong(songMetaData) //? save the song's metadata
}
func NewSongUsecase(songRepo_ repository.ISongRepo, redisRepo_ repository.IRedisSearchRepo) ISongUsecase {
	return &SongUsecase{
		songRepo:        songRepo_,
		redisSearchRepo: redisRepo_,
	}
}
