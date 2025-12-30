package http

import (
	"github.com/gin-gonic/gin"
)

type IUploadController interface {
	UploadFileToArchive(ctx *gin.Context)
}
type ISearchController interface {
	SearchSongs(ctx *gin.Context)
}

type IController interface{
	IUploadController
	ISearchController
	GetAllSongs(ctx *gin.Context)
}