package domian

import (
	"time"

	"github.com/google/uuid"
)

// PlaylistItem represents a mapping between a playlist and a song,
// including the position and when it was added.
type PlaylistItem struct {
	ID         uuid.UUID `json:"id" db:"id" gorm:"type:uuid;primary_key"`
	PlaylistID uuid.UUID `json:"playlist_id" db:"playlist_id" gorm:"type:uuid;not null"`
	SongID     uuid.UUID `json:"song_id" db:"song_id" gorm:"type:uuid;not null"`
	Position   int       `json:"position" db:"position"`
	AddedAt    time.Time `json:"added_at" db:"added_at" gorm:"type:timestamptz"`
}

// Playlist represents a user's playlist metadata.
type Playlist struct {
	ID        uuid.UUID `json:"id" db:"id" gorm:"type:uuid;primary_key"`
	UserID    uuid.UUID `json:"user_id" db:"user_id" gorm:"type:uuid;not null"`
	Name      string    `json:"name" db:"name"`
	IsPublic  bool      `json:"is_public" db:"is_public"`
	CreatedAt time.Time `json:"created_at" db:"created_at" gorm:"type:timestamptz"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" gorm:"type:timestamptz"`
}

type PlaylistEvent struct {
	ID        uuid.UUID `json:"id" db:"id" gorm:"type:uuid;primary_key"`
	PlaylistID uuid.UUID `json:"playlist_id" db:"playlist_id" gorm:"type:uuid;not null"`
	SongID    uuid.UUID `json:"song_id" db:"song_id" gorm:"type:uuid;not null"`
	Position  int        `json:"position" db:"position" gorm:"not null"`
	UserID    uuid.UUID `json:"user_id" db:"user_id" gorm:"type:uuid;not null"`
	CreatedAt time.Time `json:"created_at" db:"created_at" gorm:"type:timestamptz"`
}