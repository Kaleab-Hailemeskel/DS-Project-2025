package main

import (
	"streaming-service/config"
	"streaming-service/internal/delivery/http"
	"streaming-service/internal/usecase"
)

func main() {
	config.InitEnv()
	// Initialize usecase and controller
	songUsecase := usecase.NewStreamUsecase(config.SONG_ARCHIVE_DIR)
	songController := http.NewStreamController(songUsecase)

	// Further setup like starting the server would go here
	songServer := http.InitRouter(songController)
	songServer.Run(":" + config.SERVER_PORT) // Start the server
}
