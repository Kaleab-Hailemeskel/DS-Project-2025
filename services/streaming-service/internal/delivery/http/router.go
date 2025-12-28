package http

import (
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	// Swagger UI endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router
}
