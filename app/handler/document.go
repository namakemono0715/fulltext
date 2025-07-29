package handler

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"fulltext/search"
)

// Document 検索対象のドキュメント構造体
type Document struct {
	ID           string `json:"id"`
	ProjectCode  string `json:"project_code"`
	DocumentType string `json:"document_type"`
	TenantCode   string `json:"tenant_code"`
	Title        string `json:"title"`
	Body         string `json:"body"`
}

// IndexDocumentHandler ドキュメントをインデックスに追加するハンドラー
func IndexDocumentHandler(c *gin.Context) {
	log.Println("=== ドキュメントインデックス処理開始 ===")
	
	// URLパラメータから値を取得
	tenantCode := c.Param("tenant_code")
	projectCode := c.Param("project_code")
	documentType := c.Param("document_type")
	
	// バリデーション
	if tenantCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "テナントコードが必要です"})
		return
	}
	if projectCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "プロジェクトコードが必要です"})
		return
	}
	if documentType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ドキュメントタイプが必要です"})
		return
	}

	var doc Document
	if err := c.ShouldBindJSON(&doc); err != nil {
		log.Printf("JSONバインドエラー: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "不正なJSONフォーマット: " + err.Error()})
		return
	}

	// ドキュメントのバリデーション
	if strings.TrimSpace(doc.ID) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ドキュメントIDが必要です"})
		return
	}
	if strings.TrimSpace(doc.Title) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "タイトルが必要です"})
		return
	}

	// ドキュメントの値を設定
	doc.ID = tenantCode + ":" + projectCode + ":" + documentType + ":" + doc.ID
	doc.TenantCode = tenantCode
	doc.ProjectCode = projectCode
	doc.DocumentType = documentType

	log.Printf("インデックス処理開始 - ID: %s, タイトル: %s", doc.ID, doc.Title)

	if err := search.IndexDocument(tenantCode, doc.ID, doc); err != nil {
		log.Printf("インデックス処理エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "インデックス処理に失敗しました: " + err.Error()})
		return
	}

	log.Printf("インデックス処理完了 - ID: %s", doc.ID)
	c.JSON(http.StatusOK, gin.H{
		"status":      "indexed",
		"document_id": doc.ID,
		"tenant":      tenantCode,
	})
}

// SearchDocumentsHandler ドキュメントを検索するハンドラー
func SearchDocumentsHandler(c *gin.Context) {
	log.Println("=== ドキュメント検索処理開始 ===")
	
	tenantCode := c.Param("tenant_code")
	query := c.Query("q")

	// バリデーション
	if tenantCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "テナントコードが必要です"})
		return
	}
	if strings.TrimSpace(query) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "検索クエリが必要です"})
		return
	}

	log.Printf("検索処理開始 - テナント: %s, クエリ: %s", tenantCode, query)

	results, err := search.SearchDocuments(tenantCode, query)
	if err != nil {
		log.Printf("検索処理エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "検索処理に失敗しました: " + err.Error()})
		return
	}

	log.Printf("検索処理完了 - ヒット数: %d", len(results.Hits))
	c.JSON(http.StatusOK, gin.H{
		"results": results,
		"query":   query,
		"tenant":  tenantCode,
	})
}

// FuzzySearchDocumentsHandler ファジー検索（タイポ対応）を実行するハンドラー
func FuzzySearchDocumentsHandler(c *gin.Context) {
	log.Println("=== ファジー検索処理開始 ===")
	
	tenantCode := c.Param("tenant_code")
	query := c.Query("q")
	fuzzinessStr := c.DefaultQuery("fuzziness", "1") // デフォルトは1

	// バリデーション
	if tenantCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "テナントコードが必要です"})
		return
	}
	if strings.TrimSpace(query) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "検索クエリが必要です"})
		return
	}

	// fuzzinessの変換とバリデーション
	fuzziness := 1 // デフォルト値
	if f, err := strconv.Atoi(fuzzinessStr); err == nil && f >= 0 && f <= 2 {
		fuzziness = f
	} else if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fuzzinessは0-2の整数で指定してください"})
		return
	}

	log.Printf("ファジー検索処理開始 - テナント: %s, クエリ: %s, fuzziness: %d", tenantCode, query, fuzziness)

	results, err := search.FuzzySearchDocuments(tenantCode, query, fuzziness)
	if err != nil {
		log.Printf("ファジー検索処理エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ファジー検索処理に失敗しました: " + err.Error()})
		return
	}

	log.Printf("ファジー検索処理完了 - ヒット数: %d", len(results.Hits))
	c.JSON(http.StatusOK, gin.H{
		"results":   results,
		"query":     query,
		"fuzziness": fuzziness,
		"tenant":    tenantCode,
	})
}