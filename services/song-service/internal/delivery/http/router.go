package http

import (
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/gin-gonic/gin"
)

func RegisterUploadRoutes(router *gin.Engine, uploadController IUploadController) {
	uploadGroup := router.Group("/upload")
	{
		uploadGroup.POST("/song", uploadController.UploadFileToArchive)
	}
}
func RegisterSearchRoutes(router *gin.Engine, searchController ISearchController) {
	searchGroup := router.Group("/search")
	{
		searchGroup.GET("/songs", searchController.SearchSongs)
	}
}	

func InitRouter(uploadController IUploadController, searchController ISearchController) *gin.Engine {
	router := gin.Default()

	// Register upload routes
	RegisterUploadRoutes(router, uploadController)	
	// Register search routes
	RegisterSearchRoutes(router, searchController)
	// Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router
}