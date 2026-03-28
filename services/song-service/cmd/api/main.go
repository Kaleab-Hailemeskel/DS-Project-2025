package main

import (
	"song-service/api/config"
	"song-service/api/internal/delivery/http"
	"song-service/api/internal/repository"
	"song-service/api/internal/usecase"
	"song-service/api/pkg/media"
	"song-service/api/internal/middleware"
)

func main() {
	// Initialize main PostgreSQL repository and passing the gorm to NewSongRepository
	config.InitEnv()
	postresDb := repository.InitPostgresDB()
	redisDb := repository.InitRedisClient()
	songRepo := repository.NewSongRepository(postresDb)
	redisSearchRepo := repository.NewRedisRepository(redisDb)
	songUsecase := usecase.NewSongUsecase(songRepo, redisSearchRepo)
	mediaProcessor := media.NewMediaProcessor()
	songController := http.NewController(songUsecase, mediaProcessor)
	middleware_ := middleware.NewMiddleware(redisSearchRepo)
	// Further setup like starting the server would go here
	songServer := http.InitRouter(songController, middleware_)
	songServer.Run(":" + config.SERVER_PORT) // Start the server
}
