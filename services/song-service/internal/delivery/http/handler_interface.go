package http

import (
	"github.com/gin-gonic/gin"
)

type IUploadController interface {
	UploadFileToArchive(ctx *gin.Context)
}
