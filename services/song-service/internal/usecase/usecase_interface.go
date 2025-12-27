package usecase

import (
	"song-service/api/internal/domain"
)

type IUploadUsecase interface {
	SaveSongMetaData(songMetaData *domain.Song) (*domain.Song, error)
}

type ISearchEngineUsecase interface {
	// 1. Primary Search Logic (The core predictive suggestion feature)
	SearchSongsByTitlePrefix(titlePrefix string, pageNumber, pageLimit string) ([]*domain.Song, error)

	// 2. Search & Filtering Logic (Including Genre/Year)
	FilterSongs(query string, filters map[string]interface{}) ([]*domain.Song, error) //! not imped yet

	// 3. Data Synchronization
	IndexSong(song *domain.Song) error
	DeindexSong(song *domain.Song) error
}
