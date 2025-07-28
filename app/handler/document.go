package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"fulltext/search"
)

type Document struct {
	ID           string `json:"id"`
	ProjectCode  string `json:"project_code"`
	DocumentType string `json:"document_type"`
	TenantCode   string `json:"tenant_code"`
	Title        string `json:"title"`
	Body         string `json:"body"`
}

func IndexDocumentHandler(c *gin.Context) {
	tenant_code := c.Param("tenant_code")
	project_code := c.Param("project_code")
	document_type := c.Param("document_type")
	tenant := tenant_code
	// tenant := tenant_code + "_" + project_code + "_" + document_type

	var doc Document
	if err := c.ShouldBindJSON(&doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := search.IndexDocument(tenant, doc.ID, doc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "indexed"})
}

func SearchDocumentsHandler(c *gin.Context) {
	tenant := c.Param("tenant")
	query := c.Query("q")

	results, err := search.SearchDocuments(tenant, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}