package usecase

import (
	"context"
	"song-service/api/config"
	"song-service/api/internal/domain"
	"song-service/api/internal/repository"
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
func (s *SearchEngineUsecase) SearchSongsByTitlePrefix(titlePrefix string, offset int64, page int64) ([]*domain.Song, error) {
	newOffSet := offset + (page-1)*config.MAX_PAGE_SIZE
	return s.redisSearchRepo.SearchSongsByTitlePrefix(context.Background(), titlePrefix, newOffSet, page)
}

func NewSearchEngineUsecase(songRepo_ repository.ISongRepo) ISearchEngineUsecase {
	return &SearchEngineUsecase{
		songRepo: songRepo_,
	}
}
