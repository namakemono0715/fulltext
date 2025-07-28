package router

import (
	"github.com/gin-gonic/gin"
	"fulltext/handler"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", handler.Root)
	r.GET("/ping", handler.Ping)
	r.GET("/users", handler.Users)

	return r
}