package usecase

import (
	"context"
	"log"
	"song-service/api/internal/domain"
	"song-service/api/internal/repository"
	"strconv"
)

type SearchEngineUsecase struct {
	songRepo        repository.ISongRepo
	redisSearchRepo repository.IRedisSearchRepo
}

// DeindexSong implements ISearchEngineUsecase.
func (s *SearchEngineUsecase) DeindexSong(song *domain.Song) error {
	return s.redisSearchRepo.DeindexSong(context.Background(), song)
}

// FilterSongs implements ISearchEngineUsecase.
func (s *SearchEngineUsecase) FilterSongs(query string, filters map[string]interface{}) ([]*domain.Song, error) {
	panic("unimplemented")
}

// IndexSong implements ISearchEngineUsecase.
func (s *SearchEngineUsecase) IndexSong(song *domain.Song) error {
	return s.redisSearchRepo.IndexSong(context.Background(), song)
}

// SearchSongsByTitlePrefix implements ISearchEngineUsecase.
func (s *SearchEngineUsecase) SearchSongsByTitlePrefix(titlePrefix, pageNumber, pageLimit string) ([]*domain.Song, error) {
	// Convert offset and pageLimit to integers
	pageNumberInt, err := strconv.ParseInt(pageNumber, 10, 64)
	if err != nil {
		return nil, err
	}
	pageLimitInt, err := strconv.ParseInt(pageLimit, 10, 64)
	if err != nil {
		return nil, err
	}

	searchCache, err := s.redisSearchRepo.SearchSongsByTitlePrefix(context.Background(), titlePrefix, pageNumberInt, pageLimitInt)
	if err != nil {
		return nil, err
	}
	if searchCache == nil {
		// If no results in cache, search in the song repository (database)

		res, err := s.songRepo.GetSongByTitle(titlePrefix)
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

func NewSearchEngineUsecase(songRepo_ repository.ISongRepo, redisRepo_ repository.IRedisSearchRepo) ISearchEngineUsecase {
	return &SearchEngineUsecase{
		songRepo:        songRepo_,
		redisSearchRepo: redisRepo_,
	}
}
