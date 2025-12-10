package domain

import (
	"time"

	"github.com/google/uuid"
)

// Song represents a music track stored by the song service.
type Song struct {
	ID          uuid.UUID `json:"id" db:"id" gorm:"type:uuid;primary_key"`
	Title       string    `json:"title" db:"title"`
	Artist      string    `json:"artist" db:"artist"`
	Album       string    `json:"album" db:"album"`
	DurationSec int       `json:"duration_sec" db:"duration_sec"`
	ReleaseYear int       `json:"release_year" db:"release_year"`
	CreatedAt   time.Time `json:"created_at" db:"created_at" gorm:"type:timestamptz"`
	Genre       string    `json:"genre" db:"genre"`
}

// SongVariant represents a specific encoded variant of a Song (audio/video rendition).
type SongVariant struct {
	ID          uuid.UUID `json:"id" db:"id" gorm:"type:uuid;primary_key"`
	SongID      uuid.UUID `json:"song_id" db:"song_id" gorm:"type:uuid;not null"`
	Codec       string    `json:"codec" db:"codec"` // e.g. "aac", "opus"
	BitrateKbps int       `json:"bitrate_kbps" db:"bitrate_kbps"`
	Resolution  string    `json:"resolution,omitempty" db:"resolution"` // optional, for video
	ManifestURL string    `json:"manifest_url" db:"manifest_url"`       // points to HLS/DASH manifest
	CreatedAt   time.Time `json:"created_at" db:"created_at" gorm:"type:timestamptz"`
}
// SongEvent represents a song event
type SongEvent struct {
	ID        uuid.UUID `json:"id" db:"id" gorm:"type:uuid;primary_key"`
	SongID    uuid.UUID `json:"song_id" db:"song_id" gorm:"type:uuid;not null"`
	Title     string    `json:"title" db:"title"`
	Artist    string    `json:"artist" db:"artist"`
	DurationSec int       `json:"duration_sec" db:"duration_sec"`
	Variants  []SongVariant `json:"variants" db:"variants"`
	CreatedAt time.Time `json:"created_at" db:"created_at" gorm:"type:timestamptz"`
}

type ISongRepo interface{
	GetSong(id uuid.UUID) (*Song, error)
	SaveSong(song *Song) (*Song, error)
	GetAllSongs(mulicListPerPage, pageNumber int) ([]*Song, error)
	GetSongByArtist(artist string) ([]*Song, error)
	GetSongByTitle(title string) ([]*Song, error)
	GetSongByAlbum(album string) ([]*Song, error)
	GetSongByGenre(genre string) ([]*Song, error)
	UpdateSong(song *Song) (*Song, error)
	DeleteSong(id uuid.UUID) error
}