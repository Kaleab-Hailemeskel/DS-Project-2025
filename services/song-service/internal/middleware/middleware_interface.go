package middleware

import "github.com/gin-gonic/gin"

type IMiddleware interface {
	AuthUser(ctx *gin.Context)
}