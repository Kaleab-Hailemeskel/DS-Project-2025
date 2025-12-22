package http

import (
	"github.com/gin-gonic/gin"
)

func RegisterStreamRoutes(router *gin.Engine, streamController IStreamController) {
	streamGroup := router.Group("/streams")
	{
		streamGroup.GET("/:filename/index.m3u8", streamController.GetManifestFile)
		streamGroup.GET("/:filename/:segment", streamController.GetStreamFile)
	}
}

func InitRouter(streamController IStreamController) *gin.Engine {
	router := gin.Default()
	RegisterStreamRoutes(router, streamController)
	return router
}
