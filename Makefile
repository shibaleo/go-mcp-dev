.PHONY: dev build test health clean

APP_NAME := go-mcp-dev

# ローカル開発
dev:
	docker-compose up

# ビルド確認
build:
	docker-compose run --rm app go build ./...

# テスト
test:
	docker-compose run --rm app go test ./...

# ヘルスチェック
health:
	curl -s http://localhost:8080/health

# クリーンアップ
clean:
	docker-compose down --rmi local
