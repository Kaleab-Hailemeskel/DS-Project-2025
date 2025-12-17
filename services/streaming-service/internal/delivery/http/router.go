package http

import (
	"github.com/gin-gonic/gin"
)

func RegisterStreamRoutes(router *gin.Engine, streamController IStreamController) {
	streamGroup := router.Group("/streams")
	{
		streamGroup.GET("/get-manifest/:filename", streamController.GetManifestFile)
		streamGroup.GET("/:filename/:segment", streamController.GetStreamFile)
	}
}

func InitRouter(streamController IStreamController) *gin.Engine {
	router := gin.Default()
	RegisterStreamRoutes(router, streamController)
	return router
}
