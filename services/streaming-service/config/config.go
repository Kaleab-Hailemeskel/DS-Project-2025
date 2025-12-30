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
	REDIS_DB          int
	REDIS_PASSWORD    string
	REDIS_ADDR        string
	USER_SERVICE_PORT string
	USER_SERVICE_URL  string
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
	REDIS_DB = 0
	if value, err := strconv.Atoi(getEnv("REDIS_DB")); err == nil{
		REDIS_DB = value
	}
	REDIS_PASSWORD = getEnv("REDIS_PASSWORD")
	REDIS_ADDR = getEnv("REDIS_ADDR")
	USER_SERVICE_PORT = getEnv("USER_SERVICE_PORT")
	USER_SERVICE_URL = getEnv("USER_SERVICE_URL")
}

func getEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Environment variable %s is not set", key)
	}
	return val
}
