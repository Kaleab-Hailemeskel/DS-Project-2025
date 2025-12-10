package main

import (
	"song-service/api/config"
	"song-service/api/internal/delivery/http"
	"song-service/api/internal/repository"
	"song-service/api/internal/usecase"
)

func main() {
	// Initialize main PostgreSQL repository and passing the gorm to NewSongRepository
	postresDb := repository.InitPostgresDB()
	songRepo := repository.NewSongRepository(postresDb)
	songUsecase := usecase.NewUploadUsecase(songRepo)
	songController := http.NewUploadController(songUsecase)
	
	// Further setup like starting the server would go here
	songServer := http.InitRouter(songController)
	songServer.Run(":" + config.SERVER_PORT) // Start the server
}
