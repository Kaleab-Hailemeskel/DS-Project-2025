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
	MAX_PAGE_SIZE     int64
	REDIS_ADDR        string
	REDIS_PASSWORD    string
	REDIS_DB          int
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
	MAX_PAGE_SIZE, _ = strconv.ParseInt(getEnv("MAX_PAGE_SIZE"), 10, 64)
	if MAX_PAGE_SIZE == 0 {
		MAX_PAGE_SIZE = 10 // default max page size
	}
	REDIS_ADDR = getEnv("REDIS_ADDR")
	REDIS_PASSWORD = getEnv("REDIS_PASSWORD")
	REDIS_DB, _ = strconv.Atoi(getEnv("REDIS_DB"))
	// print loaded config for verification
	log.Printf("Config loaded: \nSONG_ARCHIVE_DIR=%s\n SERVER_PORT=%s\n POSTGRES_HOST=%s\n POSTGRES_PORT=%s\n POSTGRES_USER=%s\n POSTGRES_DB=%s\n SEGMENT_DURATION=%d\n MAX_PAGE_SIZE=%d\n REDIS_ADDR=%s\n REDIS_DB=%d",
		SONG_ARCHIVE_DIR, SERVER_PORT, POSTGRES_HOST, POSTGRES_PORT, POSTGRES_USER, POSTGRES_DB, SEGMENT_DURATION, MAX_PAGE_SIZE, REDIS_ADDR, REDIS_DB)
}

func getEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Environment variable %s is not set", key)
	}
	return val
}
