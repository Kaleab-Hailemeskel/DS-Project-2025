package main

import (
	"streaming-service/internal/usecase"
	"streaming-service/config"
	"streaming-service/internal/delivery/http"
)

func main() {
	
	// Initialize usecase and controller
	songUsecase := usecase.NewStreamUsecase(config.SONG_ARCHIVE_DIR)
	songController := http.NewStreamController(songUsecase)

	// Further setup like starting the server would go here
	songServer := http.InitRouter(songController)
	songServer.Run(":" + config.SERVER_PORT) // Start the server
}
