package search

import (
	"fmt"
	"path/filepath"
	"sync"
	"github.com/blevesearch/bleve/v2"
)

var (
	indexes   = make(map[string]bleve.Index)
	indexLock sync.Mutex
	indexPath = "./indexes" // インデックスを保存するディレクトリ
)

// getOrCreateIndex returns an existing index for the tenant or creates a new one in memory
func getOrCreateIndex(tenant string) (bleve.Index, error) {
	indexLock.Lock()
	defer indexLock.Unlock()

	if idx, exists := indexes[tenant]; exists {
		return idx, nil
	}

	// インデックスのパスを生成
	indexFilePath := filepath.Join(indexPath, fmt.Sprintf("%s.bleve", tenant))

	// 既存のインデックスを開くか、新規作成
	var idx bleve.Index
	var err error
	if bleveIndexExists(indexFilePath) {
			idx, err = bleve.Open(indexFilePath)
	} else {
			mapping := bleve.NewIndexMapping()
			idx, err = bleve.New(indexFilePath, mapping)
	}
	if err != nil {
			return nil, err
	}

	indexes[tenant] = idx
	return idx, nil
}

// IndexDocument adds a document to the tenant's index
func IndexDocument(tenant, docID string, doc interface{}) error {
	idx, err := getOrCreateIndex(tenant)
	if err != nil {
		return err
	}
	return idx.Index(docID, doc)
}

// SearchDocuments runs a simple query against the tenant's index
func SearchDocuments(tenant, query string) (*bleve.SearchResult, error) {
	idx, err := getOrCreateIndex(tenant)
	if err != nil {
		return nil, err
	}

	q := bleve.NewQueryStringQuery(query)
	req := bleve.NewSearchRequest(q)
	return idx.Search(req)
}

// bleveIndexExists checks if a Bleve index exists at the given path
func bleveIndexExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}