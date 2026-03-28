package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	// Migrate schema
	if err := ensureSchema(db); err != nil {
		log.Fatalf("migration failed: %v", err)
	}
	return &Store{db: db}
}

// ensureSchema creates the necessary tables if they do not exist.
func ensureSchema(db *sql.DB) error {
	log.Println("[Schema] Checking if tables exist...")
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}
	log.Println("[Schema] Database tables are ready.")
	return nil
}
func (s *Store) CreateUser(ctx context.Context, username, email, passwordHash string) (int, time.Time, error) {
	var id int
	var created time.Time
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	err := s.db.QueryRowContext(ctx,
		"INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id, created_at",
		username, email, passwordHash).Scan(&id, &created)
	return id, created, err
}

func (s *Store) GetUserByID(ctx context.Context, id int) (*User, error) {
	var u User
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	err := s.db.QueryRowContext(ctx,
		"SELECT id, username, email, password_hash, created_at FROM users WHERE id = $1",
		id).Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Store) GetUserByUsernameOrEmail(ctx context.Context, ident string) (*User, error) {
	var u User
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	err := s.db.QueryRowContext(ctx,
		"SELECT id, username, email, password_hash, created_at FROM users WHERE username=$1 OR email=$1 LIMIT 1",
		ident).Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
