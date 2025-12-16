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
	return router
}