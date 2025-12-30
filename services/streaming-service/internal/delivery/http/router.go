package http

import (
	"streaming-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterStreamRoutes(router *gin.RouterGroup, streamController IStreamController) {

	router.GET("/:filename/index.m3u8", streamController.GetManifestFile)
	router.GET("/:filename/:segment", streamController.GetStreamFile)

}

func InitRouter(streamController IStreamController, middleware middleware.IMiddleware) *gin.Engine {
	router := gin.Default()
	streamGroup := router.Group("/streams", middleware.AuthUser)
	RegisterStreamRoutes(streamGroup, streamController)
	return router
}
