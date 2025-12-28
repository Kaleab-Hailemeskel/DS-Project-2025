package main

// @title Song Service API
// @version 1.0
// @description API for uploading songs, segmenting to HLS, and searching metadata.
// @BasePath /
// @schemes http

import (
	"song-service/api/config"
	"song-service/api/internal/delivery/http"
	"song-service/api/internal/repository"
	"song-service/api/internal/usecase"
	"song-service/api/pkg/media"
	_ "song-service/api/docs"
)

func main() {
	// Initialize main PostgreSQL repository and passing the gorm to NewSongRepository
	config.InitEnv()
	postresDb := repository.InitPostgresDB()
	redisDb := repository.InitRedisClient()
	songRepo := repository.NewSongRepository(postresDb)
	redisSearchRepo := repository.NewRedisRepository(redisDb)
	songUsecase := usecase.NewUploadUsecase(songRepo)
	chunkerUseCase := media.NewHLSSegmenter()
	songController := http.NewUploadController(songUsecase, chunkerUseCase)
	searchUsecase := usecase.NewSearchEngineUsecase(songRepo, redisSearchRepo)
	searchController := http.NewSearchController(searchUsecase)

	// Further setup like starting the server would go here
	songServer := http.InitRouter(songController, searchController)
	songServer.Run(":" + config.SERVER_PORT) // Start the server
}
