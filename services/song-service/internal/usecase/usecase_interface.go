package usecase

import (
	"song-service/api/internal/domain"
)

type IUploadUsecase interface {
	SaveSongMetaData(songMetaData *domain.Song) (*domain.Song, error)
}

type ISearchEngineUsecase interface {
	// 1. Primary Search Logic (The core predictive suggestion feature)
	// This function remains focused on prefix matching for the main query field (Title).
	SearchSongsByTitlePrefix(titlePrefix string, offset, page int64) ([]*domain.Song, error)

	// 2. Search & Filtering Logic (Including Genre/Year)

	// Generalized search function that allows filtering by multiple fields.
	// This is useful for building a faceted search or a "search filter" UI.
	// Example: query="Bohemian", filter={"Genre": "Rock", "Year": 1975}
	FilterSongs(query string, filters map[string]interface{}) ([]*domain.Song, error)

	// // Specialized searches can be added for common use cases:
	// SearchByGenre(genre string, query string) ([]*domain.Song, error)

	// 3. Data Synchronization
	IndexSong(song *domain.Song) error
	DeindexSong(song *domain.Song) error
}
