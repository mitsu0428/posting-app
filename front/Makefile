# 掲示板アプリ フロントエンド Makefile

# 変数定義
SCHEMA_PATH = ../api/schema.yaml
OUTPUT_DIR = src/generated

.PHONY: help install build dev clean test lint typecheck codegen

# デフォルトターゲット
help:
	@echo "Available commands:"
	@echo "  install     - Install dependencies"
	@echo "  build       - Build production bundle"
	@echo "  dev         - Start development server"
	@echo "  clean       - Clean build artifacts"
	@echo "  test        - Run tests"
	@echo "  lint        - Run ESLint"
	@echo "  typecheck   - Run TypeScript type checking"
	@echo "  codegen     - Generate API client from OpenAPI schema"

# 依存関係インストール
install:
	npm install

# 本番ビルド
build:
	npm run build

# 開発サーバー起動
dev:
	npm start

# クリーンアップ
clean:
	rm -rf build
	rm -rf styled-system
	rm -rf $(OUTPUT_DIR)
	rm -rf node_modules/.cache
	rm -rf .eslintcache
	npm cache clean --force 2>/dev/null || true

# テスト実行
test:
	npm test

# ESLint実行
lint:
	npm run lint

# TypeScript型チェック
typecheck:
	npm run typecheck

# APIクライアント自動生成 (orval使用)
codegen:
	@echo "Installing orval if not present..."
	@which orval || npm install -g orval
	@echo "Generating API client from OpenAPI schema..."
	@mkdir -p $(OUTPUT_DIR)
	orval --config orval.config.js

# Docker関連
docker-build:
	docker build -t posting-app-frontend .

docker-run:
	docker run -p 3000:3000 posting-app-frontend

# プリコミットチェック
pre-commit: lint typecheck test
	@echo "All checks passed!"

# ワッチモード開発
watch:
	npm run start

# 全てのチェックを実行
check-all: install lint typecheck test build
	@echo "All checks completed successfully!"