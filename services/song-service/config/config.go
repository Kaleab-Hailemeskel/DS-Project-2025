package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	SONG_ARCHIVE_DIR  string
	SERVER_PORT       string
	POSTGRES_HOST     string
	POSTGRES_PORT     string
	POSTGRES_USER     string
	POSTGRES_PASSWORD string
	POSTGRES_DB       string
	SEGMENT_DURATION  int
)

func InitEnv() {
	// Try to load .env (for local dev), but don't fail if not found
	_ = godotenv.Load()
	SONG_ARCHIVE_DIR = getEnv("SONG_ARCHIVE_DIR")
	POSTGRES_HOST = getEnv("POSTGRES_HOST")
	POSTGRES_PORT = getEnv("POSTGRES_PORT")
	POSTGRES_USER = getEnv("POSTGRES_USER")
	POSTGRES_PASSWORD = getEnv("POSTGRES_PASSWORD")
	POSTGRES_DB = getEnv("POSTGRES_DB")
	// Render provides SERVER_PORT automatically
	SERVER_PORT = os.Getenv("SERVER_PORT")
	if SERVER_PORT == "" {
		SERVER_PORT = "8080" // fallback for local dev
	}
	SEGMENT_DURATION, _ = strconv.Atoi(getEnv("SEGMENT_DURATION"))
	if SEGMENT_DURATION == 0 {
		SEGMENT_DURATION = 10 // default segment duration
	}
}

func getEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Environment variable %s is not set", key)
	}
	return val
}
