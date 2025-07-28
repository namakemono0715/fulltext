package router

import (
	"fulltext/handler"
	"fulltext/middleware"
	"github.com/gin-gonic/gin"
	"os"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	apiKey := os.Getenv("API_KEY")

	auth := middleware.AuthMiddleware(apiKey)

	secured := r.Group("/:tenant_code/:project_code/:document_type", auth)
	{
		secured.POST("/documents", handler.IndexDocumentHandler)
		secured.GET("/search", handler.SearchDocumentsHandler)
	}

	return r
}