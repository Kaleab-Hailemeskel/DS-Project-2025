package http

import (
	"github.com/gin-gonic/gin"
	"song-service/api/internal/middleware"
)

//	func RegisterUploadRoutes(router *gin.Engine, uploadController IUploadController) {
//		router.POST("/upload", uploadController.UploadFileToArchive)
//	}
//
//	func RegisterSearchRoutes(router *gin.Engine, searchController ISearchController) {
//		router.GET("/search", searchController.SearchSongs)
//	}
func RegisterUserRoutes(router *gin.RouterGroup, controller IController) {
	router.GET("/search", controller.SearchSongs)
	router.POST("/upload", controller.UploadFileToArchive)
	router.GET("/get-all", controller.GetAllSongs)
}


func InitRouter(controller IController, middleWare middleware.IMiddleware) *gin.Engine {
	router := gin.Default()

	// Register upload routes, search routes
	authedRoute := router.Group("", middleWare.AuthUser)
	RegisterUserRoutes(authedRoute, controller)
	return router
}
