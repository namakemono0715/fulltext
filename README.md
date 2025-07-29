# Fulltext App – Docker Dev Template

このプロジェクトは、Go（Gin）をベースにした 全文検索APIアプリケーションです。
※最低限の機能しか実装されていません。  
Docker と [Air](https://github.com/air-verse/air) を使って、ソースコードの変更を即時に反映できるホットリロード開発が可能です。

## 📦 構成技術

- Go 1.24
- Gin Web Framework
- Air（ホットリロード）
- Docker / Docker Compose
- Caddy（リバースプロキシ）
- Makefile（開発補助コマンド）

---

## 🚀 セットアップ

### 1. ホストの `/etc/hosts` に追記

```bash
127.0.0.1   fulltext.localhos
```

> ※ `sudo vi /etc/hosts` 等で編集してください。

---

## 🔧 開発コマンド一覧（Makefile）

| コマンド        | 説明                                 |
|-----------------|--------------------------------------|
| `make dev`      | ホットリロード付きで開発起動         |
| `make build`    | ビルド付きで起動（Dockerfile変更時） |
| `make down`     | 停止                                 |
| `make restart`  | 再起動（down → up）                  |
| `make shell`    | app コンテナに入る                   |
| `make logs`     | リアルタイムログを見る               |
| `make air`      | app コンテナ内で air を再実行        |

---

## 🌐 アクセス

- アプリ： [http://fulltext.localhos](http://fulltext.localhos)

※ Caddy がリバースプロキシしており、ポート番号不要でアクセスできます。  
※ 既に `puma-dev` などが `80` 番を使っている場合は、Caddy のポートを変更して `.local:8081` などにしてください。

---

## 📋 API 仕様

### 認証

すべてのAPIエンドポイントには、環境変数 `API_KEY` で設定されたAPIキーによる認証が必要です。

**Header**
```
Authorization: Bearer YOUR_API_KEY
```

### ベースURL

```
http://fulltext.localhos/{tenant_code}/{project_code}/{document_type}
```

**パラメータ説明**
- `tenant_code`: テナント識別子
- `project_code`: プロジェクト識別子  
- `document_type`: ドキュメントタイプ

### エンドポイント

#### 1. ドキュメントのインデックス化

**POST** `/{tenant_code}/{project_code}/{document_type}/documents`

ドキュメントを検索インデックスに追加します。

**リクエストボディ**
```json
{
  "id": "document_001",
  "title": "サンプルドキュメント",
  "body": "これはサンプルのドキュメント本文です。"
}
```

**レスポンス（成功）**
```json
{
  "status": "indexed",
  "document_id": "tenant1:project1:manual:document_001",
  "tenant": "tenant1"
}
```

**レスポンス（エラー）**
```json
{
  "error": "ドキュメントIDが必要です"
}
```

#### 2. ドキュメント検索

**GET** `/{tenant_code}/{project_code}/{document_type}/search?q={query}`

インデックス化されたドキュメントを検索します。

**クエリパラメータ**
- `q`: 検索クエリ文字列（必須）

**レスポンス（成功）**
```json
{
  "results": {
    "hits": [
      {
        "id": "tenant1:project1:manual:document_001",
        "score": 0.8,
        "fields": {
          "title": "サンプルドキュメント",
          "body": "これはサンプルのドキュメント本文です。"
        }
      }
    ],
    "total_hits": 1,
    "max_score": 0.8
  },
  "query": "サンプル",
  "tenant": "tenant1"
}
```

**レスポンス（エラー）**
```json
{
  "error": "検索クエリが必要です"
}
```

#### 3. ファジー検索（タイポ対応）

**GET** `/{tenant_code}/{project_code}/{document_type}/fuzzy-search?q={query}&fuzziness={level}`

タイプミスに対応したファジー検索を実行します。

**クエリパラメータ**
- `q`: 検索クエリ文字列（必須）
- `fuzziness`: あいまい度レベル 0-2（省略可、デフォルト: 1）
  - `0`: 完全一致のみ
  - `1`: 1文字の違いまで許可（推奨）
  - `2`: 2文字の違いまで許可

**レスポンス（成功）**
```json
{
  "results": {
    "hits": [
      {
        "id": "tenant1:project1:manual:document_001",
        "score": 0.7,
        "fields": {
          "title": "サンプルドキュメント",
          "body": "これはサンプルのドキュメント本文です。"
        }
      }
    ],
    "total_hits": 1,
    "max_score": 0.7
  },
  "query": "サンプる",
  "fuzziness": 1,
  "tenant": "tenant1"
}
```

### 使用例

#### ドキュメントの追加
```bash
curl -X POST \
  http://fulltext.localhos/tenant1/project1/manual/documents \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "doc001",
    "title": "ユーザーマニュアル",
    "body": "このマニュアルはシステムの使用方法を説明します。"
  }'
```

#### ドキュメントの検索
```bash
curl -X GET \
  "http://fulltext.localhos/tenant1/project1/manual/search?q=ユーザー" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

#### ファジー検索（タイポ対応）
```bash
# fuzziness=1（1文字違いまで許可）
curl -X GET \
  "http://fulltext.localhos/tenant1/project1/manual/fuzzy-search?q=ユーザ&fuzziness=1" \
  -H "Authorization: Bearer YOUR_API_KEY"

# fuzziness=2（2文字違いまで許可）
curl -X GET \
  "http://fulltext.localhos/tenant1/project1/manual/fuzzy-search?q=ユザー&fuzziness=2" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

---

## 📁 ディレクトリ構成

```
.
├── app/
│   ├── main.go              # メインアプリケーション
│   ├── handler/             # HTTPハンドラー
│   │   └── document.go      # ドキュメント関連API
│   ├── router/              # ルーティング設定
│   │   └── router.go        # API ルート定義
│   ├── middleware/          # ミドルウェア
│   │   └── auth.go          # 認証ミドルウェア
│   ├── search/              # 検索エンジン
│   │   └── index.go         # Bleve検索インデックス
│   ├── indexes/             # 検索インデックスファイル（自動生成）
│   └── tmp/                 # 一時ファイル
├── indexes/                 # 検索インデックスファイル（自動生成）
├── docker-compose.yml       # Docker構成
├── Dockerfile              # Dockerイメージ定義
├── Caddyfile               # リバースプロキシ設定
├── Makefile                # 開発コマンド
└── README.md               # このファイル
```

---

## 🛠 Tips

- `Air` により `app/` 配下の `.go` ファイルを変更すると自動で再ビルド＆再起動されます。
- 本番では `CMD ["./server"]` でビルド済みバイナリを起動する形に変更してください。
- 検索インデックスは `indexes/` ディレクトリに自動生成されます（Gitで管理されません）。
- API認証には環境変数 `API_KEY` を設定してください。
- テナント毎に独立した検索インデックスが作成されます。
- **ファジー検索対応**: タイプミスに強い検索が可能です（fuzziness 0-2で調整可能）。
- 通常検索は高速、ファジー検索はタイプミス対応で使い分けてください。

---

## 🔧 環境変数

アプリケーションで使用する環境変数：

| 変数名    | 説明                 | 必須 | デフォルト値 |
|-----------|---------------------|------|-------------|
| `API_KEY` | API認証キー          | ✅   | なし        |
| `GIN_MODE`| Ginの動作モード      | ❌   | debug       |

**設定例（docker-compose.yml）**
```yaml
environment:
  - API_KEY=your-secret-api-key
  - GIN_MODE=release
```

---

## 📜 ライセンス

MIT License
