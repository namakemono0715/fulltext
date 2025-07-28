# 開発用コンテナ起動（ホットリロード対応）
dev:
	docker compose up

# ビルドして開発用コンテナ起動（Dockerfileの変更などがあった場合）
build:
	docker compose up --build

# 停止
down:
	docker compose down

# 再起動（down → up）
restart:
	docker compose down && docker compose up

# コンテナに入る（app）
shell:
	docker compose exec app /bin/bash

# ログを見る（リアルタイム）
logs:
	docker compose logs -f

# Air 再起動（コンテナ内から）
air:
	docker compose exec app air