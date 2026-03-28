package main

import (
	"streaming-service/config"
	"streaming-service/internal/delivery/http"
	"streaming-service/internal/middleware"
	"streaming-service/internal/repository"
	"streaming-service/internal/usecase"
)

func main() {
	config.InitEnv()
	// Initialize usecase and controller
	redisClient := repository.InitRedisClient()
	redisSearchRepo := repository.NewRedisRepository(redisClient)
	songUsecase := usecase.NewStreamUsecase(config.SONG_ARCHIVE_DIR)
	middleware_ := middleware.NewMiddleware(redisSearchRepo)
	songController := http.NewStreamController(songUsecase)
	// Further setup like starting the server would go here
	songServer := http.InitRouter(songController, middleware_)
	songServer.Run(":" + config.SERVER_PORT) // Start the server
}
