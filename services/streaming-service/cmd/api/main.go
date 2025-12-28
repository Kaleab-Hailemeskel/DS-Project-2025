package main

// @title Streaming Service API
// @version 1.0
// @description API for streaming HLS content (manifests and segments)
// @BasePath /
// @schemes http

import (
	"streaming-service/internal/usecase"
	"streaming-service/config"
	"streaming-service/internal/delivery/http"
	_ "streaming-service/docs"
)

func main() {
	// Initialize environment variables (SERVER_PORT, SONG_ARCHIVE_DIR)
	config.InitEnv()

	// Initialize usecase and controller
	songUsecase := usecase.NewStreamUsecase(config.SONG_ARCHIVE_DIR)
	songController := http.NewStreamController(songUsecase)

	// Start HTTP server with routes (including Swagger UI)
	songServer := http.InitRouter(songController)
	songServer.Run(":" + config.SERVER_PORT)
}
