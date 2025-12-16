package http

import (
	"github.com/gin-gonic/gin"
)

type IStreamController interface {
	GetStreamFile(ctx *gin.Context)
	GetManifasteFile(ctx *gin.Context)
}

