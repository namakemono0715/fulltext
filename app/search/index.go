package search

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"github.com/blevesearch/bleve/v2"
)

var (
	indexes   = make(map[string]bleve.Index)
	indexLock sync.Mutex
	indexPath = "./indexes" // インデックスを保存するディレクトリ
)

// getOrCreateIndex テナント用の既存インデックスを返すか、新しいインデックスを作成する
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
		if err != nil {
			return nil, fmt.Errorf("既存インデックスの読み込みに失敗しました: %w", err)
		}
	} else {
		mapping := bleve.NewIndexMapping()
		idx, err = bleve.New(indexFilePath, mapping)
		if err != nil {
			return nil, fmt.Errorf("新規インデックスの作成に失敗しました: %w", err)
		}
	}

	indexes[tenant] = idx
	return idx, nil
}

// IndexDocument テナントのインデックスにドキュメントを追加する
func IndexDocument(tenant, docID string, doc interface{}) error {
	if docID == "" {
		return fmt.Errorf("ドキュメントIDは空にできません")
	}
	if tenant == "" {
		return fmt.Errorf("テナントは空にできません")
	}
	
	idx, err := getOrCreateIndex(tenant)
	if err != nil {
		return fmt.Errorf("インデックスの取得に失敗しました: %w", err)
	}
	
	if err := idx.Index(docID, doc); err != nil {
		return fmt.Errorf("ドキュメントのインデックス化に失敗しました: %w", err)
	}
	
	return nil
}

// SearchDocuments テナントのインデックスに対してクエリを実行し、検索結果を返す
func SearchDocuments(tenant, query string) (*bleve.SearchResult, error) {
	if tenant == "" {
		return nil, fmt.Errorf("テナントは空にできません")
	}
	if query == "" {
		return nil, fmt.Errorf("検索クエリは空にできません")
	}
	
	idx, err := getOrCreateIndex(tenant)
	if err != nil {
		return nil, fmt.Errorf("インデックスの取得に失敗しました: %w", err)
	}

	// クエリを作成して検索を実行
	q := bleve.NewQueryStringQuery(query)
	req := bleve.NewSearchRequest(q)
	
	result, err := idx.Search(req)
	if err != nil {
		return nil, fmt.Errorf("検索の実行に失敗しました: %w", err)
	}
	
	return result, nil
}

// FuzzySearchDocuments ファジー検索を実行する（タイポに対応）
func FuzzySearchDocuments(tenant, query string, fuzziness int) (*bleve.SearchResult, error) {
	if tenant == "" {
		return nil, fmt.Errorf("テナントは空にできません")
	}
	if query == "" {
		return nil, fmt.Errorf("検索クエリは空にできません")
	}
	if fuzziness < 0 || fuzziness > 2 {
		return nil, fmt.Errorf("fuzzinessは0-2の範囲で指定してください")
	}
	
	idx, err := getOrCreateIndex(tenant)
	if err != nil {
		return nil, fmt.Errorf("インデックスの取得に失敗しました: %w", err)
	}

	// ファジー検索クエリを作成
	q := bleve.NewFuzzyQuery(query)
	q.SetFuzziness(fuzziness)
	req := bleve.NewSearchRequest(q)
	
	result, err := idx.Search(req)
	if err != nil {
		return nil, fmt.Errorf("ファジー検索の実行に失敗しました: %w", err)
	}
	
	return result, nil
}

// 指定されたパスにBleveインデックスが存在するかチェックする
func bleveIndexExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// CloseIndex 指定されたテナントのインデックスを閉じる
func CloseIndex(tenant string) error {
	indexLock.Lock()
	defer indexLock.Unlock()

	idx, exists := indexes[tenant]
	if !exists {
		return fmt.Errorf("テナント '%s' のインデックスが見つかりません", tenant)
	}

	if err := idx.Close(); err != nil {
		return fmt.Errorf("インデックスのクローズに失敗しました: %w", err)
	}

	delete(indexes, tenant)
	return nil
}

// CloseAllIndexes すべてのインデックスを閉じる
func CloseAllIndexes() error {
	indexLock.Lock()
	defer indexLock.Unlock()

	var errors []string
	for tenant, idx := range indexes {
		if err := idx.Close(); err != nil {
			errors = append(errors, fmt.Sprintf("テナント '%s': %v", tenant, err))
		}
	}

	// マップをクリア
	indexes = make(map[string]bleve.Index)

	if len(errors) > 0 {
		return fmt.Errorf("一部のインデックスのクローズに失敗しました: %v", errors)
	}

	return nil
}