package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	SONG_ARCHIVE_DIR string
	SERVER_PORT      string
)

func InitEnv() {
	// Try to load .env (for local dev), but don't fail if not found
	_ = godotenv.Load()
	SONG_ARCHIVE_DIR = getEnv("SONG_ARCHIVE_DIR")

	// Render provides SERVER_PORT automatically
	SERVER_PORT = os.Getenv("SERVER_PORT")
	if SERVER_PORT == "" {
		SERVER_PORT = "8080" // fallback for local dev
	}
}

func getEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Environment variable %s is not set", key)
	}
	return val
}
