package http

import (
	"github.com/gin-gonic/gin"
)

func RegisterUploadRoutes(router *gin.Engine, uploadController IUploadController) {
	router.POST("/upload", uploadController.UploadFileToArchive)
}
func RegisterSearchRoutes(router *gin.Engine, searchController ISearchController) {
	router.GET("/search", searchController.SearchSongs)
}	

func InitRouter(uploadController IUploadController, searchController ISearchController) *gin.Engine {
	router := gin.Default()

	// Register upload routes
	RegisterUploadRoutes(router, uploadController)	
	// Register search routes
	RegisterSearchRoutes(router, searchController)
	return router
}