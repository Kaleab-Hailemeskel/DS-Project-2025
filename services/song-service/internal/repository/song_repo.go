package repository

import (
	"errors"
	"fmt"
	"song-service/api/internal/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ...existing code...

// --- PostgreSQL Repository Implementation ---

// SongRepository implements ISongRepo using GORM for PostgreSQL.
// This mirrors the MongoDB structure by holding the database connection.
type SongRepository struct {
	db *gorm.DB
}

// NewSongRepository creates a new repository for the given GORM DB instance.
func NewSongRepository(db *gorm.DB) ISongRepo {
	// AutoMigrate is called here to ensure the table structure is correct.
	// In a real application, this might be handled by migration tools.
	err := db.AutoMigrate(&domain.Song{})
	if err != nil {
		// Panic is acceptable here as the application cannot run without a working DB schema.
		panic("Failed to auto-migrate Song table: " + err.Error())
	}

	return &SongRepository{db: db}
}

// GetSong retrieves a single song by its UUID.
func (r *SongRepository) GetSong(id uuid.UUID) (*domain.Song, error) {
	var song domain.Song
	// First() finds the first record that matches the condition.
	result := r.db.First(&song, "id = ?", id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("song with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to fetch song: %w", result.Error)
	}

	return &song, nil
}

// SaveSong creates a new Song record in the database.
func (r *SongRepository) SaveSong(song *domain.Song) (*domain.Song, error) {
	// Ensure ID and CreatedAt are set if they are zero values (a common practice).
	if song.ID == uuid.Nil {
		song.ID = uuid.New()
	}
	if song.CreatedAt.IsZero() {
		song.CreatedAt = time.Now().UTC()
	}

	// Create() inserts a new record.
	result := r.db.Create(song)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to save song: %w", result.Error)
	}

	return song, nil
}

// GetAllSongs retrieves a list of songs with pagination.
func (r *SongRepository) GetAllSongs(limitPerPage, pageNumber int) ([]*domain.Song, error) {
	var songs []*domain.Song

	// Calculate offset for pagination. Ensure pageNumber is at least 1.
	if pageNumber < 1 {
		pageNumber = 1
	}
	offset := (pageNumber - 1) * limitPerPage

	// Use Scopes for cleaner query building (pagination, ordering).
	result := r.db.Limit(limitPerPage).Offset(offset).Order("created_at desc").Find(&songs)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch all songs: %w", result.Error)
	}

	return songs, nil
}

// GetSongByArtist retrieves all songs by a specific artist.
func (r *SongRepository) GetSongByArtist(artist string) ([]*domain.Song, error) {
	var songs []*domain.Song
	// Where() applies a condition.
	result := r.db.Where("artist = ?", artist).Find(&songs)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch songs by artist: %w", result.Error)
	}

	return songs, nil
}

// GetSongByTitle retrieves all songs matching a title.
func (r *SongRepository) GetSongByTitle(title string) ([]*domain.Song, error) {
	var songs []*domain.Song
	// Using LIKE for partial matching (common for title searches).
	// For case-insensitivity on PostgreSQL, you might use ILIKE.
	result := r.db.Where("title LIKE ?", "%"+title+"%").Find(&songs)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch songs by title: %w", result.Error)
	}

	return songs, nil
}

// GetSongByAlbum retrieves all songs by a specific album.
func (r *SongRepository) GetSongByAlbum(album string) ([]*domain.Song, error) {
	var songs []*domain.Song
	result := r.db.Where("album = ?", album).Find(&songs)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch songs by album: %w", result.Error)
	}

	return songs, nil
}

// GetSongByGenre retrieves all songs by a specific genre.
func (r *SongRepository) GetSongByGenre(genre string) ([]*domain.Song, error) {
	var songs []*domain.Song
	result := r.db.Where("genre = ?", genre).Find(&songs)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch songs by genre: %w", result.Error)
	}

	return songs, nil
}

// UpdateSong updates an existing Song record.
func (r *SongRepository) UpdateSong(song *domain.Song) (*domain.Song, error) {
	// First, check if the record exists to ensure atomicity/error handling is clean.
	if song.ID == uuid.Nil {
		return nil, errors.New("cannot update song: ID is required")
	}

	// Save() handles updating if the primary key exists.
	// We use Omit("created_at") to prevent GORM from updating the creation timestamp.
	result := r.db.Omit("created_at").Save(song)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to update song: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("no song found or updated with ID %s", song.ID)
	}

	return song, nil
}

// DeleteSong deletes a single song by its UUID.
func (r *SongRepository) DeleteSong(id uuid.UUID) error {
	// Delete() performs a hard delete from the database.
	result := r.db.Delete(&domain.Song{}, "id = ?", id)

	if result.Error != nil {
		return fmt.Errorf("failed to delete song: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no song found with ID %s to delete", id)
	}

	return nil
}

// ...existing code...
