package usecase

import (
	"context"
	"song-service/api/config"
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
func (s *SearchEngineUsecase) SearchSongsByTitlePrefix(titlePrefix, offset, page string) ([]*domain.Song, error) {
	// Convert offset and page to integers
	offsetInt, err := strconv.ParseInt(offset, 10, 64)
	if err != nil {
		return nil, err
	}
	pageInt, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		return nil, err
	}
	// Calculate the new offset based on page and offset
	newOffSet := offsetInt + (pageInt-1)*config.MAX_PAGE_SIZE
	return s.redisSearchRepo.SearchSongsByTitlePrefix(context.Background(), titlePrefix, newOffSet, pageInt)
}

func NewSearchEngineUsecase(songRepo_ repository.ISongRepo, redisRepo_ repository.IRedisSearchRepo) ISearchEngineUsecase {
	return &SearchEngineUsecase{
		songRepo:        songRepo_,
		redisSearchRepo: redisRepo_,
	}
}
