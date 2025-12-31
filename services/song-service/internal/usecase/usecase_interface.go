package usecase

import (
	"song-service/api/internal/domain"

	"github.com/google/uuid"
)

type IUploadUsecase interface {
	SaveSongMetaData(songMetaData *domain.Song) (*domain.Song, error)
}

type ISearchEngineUsecase interface {
	// 1. Primary Search Logic (The core predictive suggestion feature)
	SearchSongsByPrefix(titlePrefix string, pageNumber, pageLimit string) ([]*domain.Song, error)

	// 2. Search & Filtering Logic (Including Genre/Year)
	FilterSongs(query string, filters map[string]interface{}) ([]*domain.Song, error) //! not imped yet

	// 3. Data Synchronization
	IndexSong(song *domain.Song) error
	DeindexSong(song *domain.Song) error
}

type ISongUsecase interface { // inheriting both interfaces
	IUploadUsecase
	ISearchEngineUsecase
	GetAllSong(pageNumber, pageLimit string) ([]*domain.Song, error)
	SaveBlobImage(id uuid.UUID, blob []byte) error
}
