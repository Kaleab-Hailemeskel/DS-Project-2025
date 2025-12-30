CREATE TABLE IF NOT EXISTS songs (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    artist VARCHAR(255) NOT NULL,
    album VARCHAR(255),
    duration_sec INTEGER NOT NULL,
    release_year INTEGER,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    genre VARCHAR(100),
    cover_art_blob BYTEA
);