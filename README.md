# Fulltext App – Docker Dev Template

このプロジェクトは、Go（Gin）をベースにした Web アプリケーションの開発テンプレートです。  
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
127.0.0.1   fulltext.local
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

- アプリ： [http://fulltext.local](http://fulltext.local)

※ Caddy がリバースプロキシしており、ポート番号不要でアクセスできます。  
※ 既に `puma-dev` などが `80` 番を使っている場合は、Caddy のポートを変更して `.local:8081` などにしてください。

---

## 📁 ディレクトリ構成（例）

```
.
├── app/
│   ├── main.go
│   ├── handler/
│   └── ...
├── docker-compose.yml
├── Dockerfile
├── Caddyfile
├── Makefile
└── README.md
```

---

## 🛠 Tips

- `Air` により `app/` 配下の `.go` ファイルを変更すると自動で再ビルド＆再起動されます。
- 本番では `CMD ["./server"]` でビルド済みバイナリを起動する形に変更してください。

---

## 📜 ライセンス

MIT License
