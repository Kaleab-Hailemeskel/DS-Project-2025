package http

import (
	"github.com/gin-gonic/gin"
)

func RegisterUploadRoutes(router *gin.Engine, uploadController IUploadController) {
	uploadGroup := router.Group("/upload")
	{
		uploadGroup.POST("/song", uploadController.UploadFileToArchive)
	}
}

func InitRouter(uploadController IUploadController) *gin.Engine {
	router := gin.Default()

	// Register upload routes
	RegisterUploadRoutes(router, uploadController)
	return router
}