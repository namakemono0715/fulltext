package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func Users(c *gin.Context) {
	users := []map[string]interface{}{
		{"id": 1, "name": "Alice"},
		{"id": 2, "name": "Bob"},
	}
	c.JSON(http.StatusOK, users)
}