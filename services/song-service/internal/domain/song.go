package domain

import (
	"time"

	"github.com/google/uuid"
)

// Song represents a music track stored by the song service.
type Song struct {
	ID           uuid.UUID `json:"id" db:"id" gorm:"type:uuid;primary_key"`
	Title        string    `json:"title" db:"title"`
	Artist       string    `json:"artist" db:"artist"`
	Album        string    `json:"album" db:"album"`
	DurationSec  int       `json:"duration_sec" db:"duration_sec"`
	ReleaseYear  int       `json:"release_year" db:"release_year"`
	CreatedAt    time.Time `json:"created_at" db:"created_at" gorm:"type:timestamptz"`
	Genre        string    `json:"genre" db:"genre"`
	CoverArtBlob []byte    `json:"cover_art_blob" db:"cover_art_blob"`
}
