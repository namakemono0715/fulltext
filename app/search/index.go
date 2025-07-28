package search

import (
	"sync"

	"github.com/blevesearch/bleve/v2"
)

var (
	indexes   = make(map[string]bleve.Index)
	indexLock sync.Mutex
)

// getOrCreateIndex returns an existing index for the tenant or creates a new one in memory
func getOrCreateIndex(tenant string) (bleve.Index, error) {
	indexLock.Lock()
	defer indexLock.Unlock()

	if idx, exists := indexes[tenant]; exists {
		return idx, nil
	}

	mapping := bleve.NewIndexMapping()
	newIndex, err := bleve.NewMemOnly(mapping)
	if err != nil {
		return nil, err
	}

	indexes[tenant] = newIndex
	return newIndex, nil
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