package http

import (
	"github.com/gin-gonic/gin"
)

type IStreamController interface {
	GetStreamFile(ctx *gin.Context)
	GetManifestFile(ctx *gin.Context)
}
