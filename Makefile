# 掲示板アプリ 開発環境 Makefile

# カラー出力設定
GREEN = \033[0;32m
YELLOW = \033[1;33m
RED = \033[0;31m
NC = \033[0m # No Color

# Docker Compose ファイル
DOCKER_COMPOSE = docker-compose.yml
DB_CONTAINER = posting-app-db
BACKEND_CONTAINER = posting-app-backend
FRONTEND_CONTAINER = posting-app-frontend

.PHONY: help dev up down clean clean-all logs db-logs backend-logs frontend-logs \
        restart restart-backend restart-frontend restart-db \
        db-shell db-migrate db-seed \
        frontend-install frontend-build frontend-test frontend-lint frontend-typecheck \
        backend-build backend-test \
        cache-clear docker-clean volume-clean \
        health status check-all reset

# デフォルトターゲット - ヘルプを表示
help:
	@echo "$(GREEN)掲示板アプリ 開発環境 Makefile$(NC)"
	@echo ""
	@echo "$(YELLOW)🚀 開発環境管理:$(NC)"
	@echo "  dev               - 開発環境を起動（DB+Backend+Frontend）"
	@echo "  up                - 全サービスをバックグラウンドで起動"
	@echo "  down              - 全サービスを停止"
	@echo "  restart           - 全サービスを再起動"
	@echo "  status            - サービスの状態を確認"
	@echo "  logs              - 全サービスのログを表示"
	@echo ""
	@echo "$(YELLOW)🗄️  データベース管理:$(NC)"
	@echo "  db-up             - データベースのみ起動"
	@echo "  db-migrate        - データベースマイグレーション実行"
	@echo "  db-seed           - シードデータ投入"
	@echo "  db-shell          - データベースに接続"
	@echo "  db-logs           - データベースログを表示"
	@echo "  restart-db        - データベースを再起動"
	@echo ""
	@echo "$(YELLOW)🔧 バックエンド管理:$(NC)"
	@echo "  backend-up        - バックエンドのみ起動"
	@echo "  backend-build     - バックエンドをビルド"
	@echo "  backend-test      - バックエンドテスト実行"
	@echo "  backend-logs      - バックエンドログを表示"
	@echo "  restart-backend   - バックエンドを再起動"
	@echo ""
	@echo "$(YELLOW)🎨 フロントエンド管理:$(NC)"
	@echo "  frontend-up       - フロントエンドのみ起動"
	@echo "  frontend-install  - フロントエンド依存関係インストール"
	@echo "  frontend-build    - フロントエンドビルド"
	@echo "  frontend-test     - フロントエンドテスト実行"
	@echo "  frontend-lint     - ESLint実行"
	@echo "  frontend-typecheck - TypeScriptチェック"
	@echo "  frontend-logs     - フロントエンドログを表示"
	@echo "  restart-frontend  - フロントエンドを再起動"
	@echo "  api-generate      - OpenAPIからAPIクライアント生成"
	@echo ""
	@echo "$(YELLOW)🧹 クリーンアップ:$(NC)"
	@echo "  clean             - ビルドキャッシュをクリア"
	@echo "  clean-all         - 全キャッシュ・ボリューム・イメージを削除"
	@echo "  cache-clear       - Node.js/Dockerキャッシュを完全削除"
	@echo "  docker-clean      - 未使用Dockerリソースを削除"
	@echo "  volume-clean      - Dockerボリュームを削除"
	@echo "  reset             - 完全リセット（clean-all + 再ビルド）"
	@echo ""
	@echo "$(YELLOW)✅ 品質チェック:$(NC)"
	@echo "  check-all         - 全品質チェック実行"
	@echo "  health            - ヘルスチェック実行"
	@echo ""

# ====================
# 🚀 開発環境管理
# ====================

# 開発環境を起動（フォアグラウンド）
dev: cache-clear
	@echo "$(GREEN)🚀 開発環境を起動中...$(NC)"
	docker-compose up --build

# 全サービスをバックグラウンドで起動
up: cache-clear
	@echo "$(GREEN)🚀 全サービスをバックグラウンドで起動中...$(NC)"
	docker-compose up -d --build
	@echo "$(GREEN)✅ 起動完了！$(NC)"
	@echo "📱 フロントエンド: http://localhost:3000"
	@echo "🔧 バックエンド: http://localhost:8080/api"
	@echo "🗄️  データベース: localhost:5432"

# 全サービスを停止
down:
	@echo "$(YELLOW)🛑 全サービスを停止中...$(NC)"
	docker-compose down

# 全サービスを再起動
restart: down up

# サービスの状態を確認
status:
	@echo "$(GREEN)📊 サービス状態:$(NC)"
	docker-compose ps

# 全サービスのログを表示
logs:
	docker-compose logs -f

# ====================
# 🗄️ データベース管理
# ====================

# データベースのみ起動
db-up:
	@echo "$(GREEN)🗄️  データベースを起動中...$(NC)"
	docker-compose up -d db

# データベースマイグレーション実行
db-migrate: db-up
	@echo "$(GREEN)🔄 データベースマイグレーション実行中...$(NC)"
	@sleep 5  # データベース起動を待つ
	docker-compose exec db psql -U postgres -d posting_app -f /docker-entrypoint-initdb.d/001_initial_schema.sql || echo "Migration already applied"

# シードデータ投入
db-seed: db-migrate
	@echo "$(GREEN)🌱 シードデータを投入中...$(NC)"
	docker-compose exec db psql -U postgres -d posting_app -f /docker-entrypoint-initdb.d/002_seed_data.sql || echo "Seed data already exists"

# データベースに接続
db-shell:
	@echo "$(GREEN)🔗 データベースに接続中...$(NC)"
	docker-compose exec db psql -U postgres -d posting_app

# データベースログを表示
db-logs:
	docker-compose logs -f db

# データベースを再起動
restart-db:
	@echo "$(YELLOW)🔄 データベースを再起動中...$(NC)"
	docker-compose restart db

# ====================
# 🔧 バックエンド管理
# ====================

# バックエンドのみ起動
backend-up: db-up
	@echo "$(GREEN)🔧 バックエンドを起動中...$(NC)"
	docker-compose up -d backend

# バックエンドをビルド
backend-build:
	@echo "$(GREEN)🏗️  バックエンドをビルド中...$(NC)"
	docker-compose build backend

# バックエンドテスト実行
backend-test:
	@echo "$(GREEN)🧪 バックエンドテストを実行中...$(NC)"
	cd backend && go test ./...

# バックエンドログを表示
backend-logs:
	docker-compose logs -f backend

# バックエンドを再起動
restart-backend:
	@echo "$(YELLOW)🔄 バックエンドを再起動中...$(NC)"
	docker-compose restart backend

# ====================
# 🎨 フロントエンド管理
# ====================

# フロントエンドのみ起動
frontend-up: backend-up
	@echo "$(GREEN)🎨 フロントエンドを起動中...$(NC)"
	docker-compose up -d frontend

# フロントエンド依存関係インストール
frontend-install:
	@echo "$(GREEN)📦 フロントエンド依存関係をインストール中...$(NC)"
	cd front && npm install --legacy-peer-deps

# フロントエンドビルド
frontend-build: frontend-install
	@echo "$(GREEN)🏗️  フロントエンドをビルド中...$(NC)"
	cd front && npm run api:generate && npm run panda && npm run build

# フロントエンドテスト実行
frontend-test:
	@echo "$(GREEN)🧪 フロントエンドテストを実行中...$(NC)"
	cd front && npm test -- --watchAll=false

# ESLint実行
frontend-lint:
	@echo "$(GREEN)🔍 ESLintを実行中...$(NC)"
	cd front && npm run lint

# TypeScriptチェック
frontend-typecheck:
	@echo "$(GREEN)📝 TypeScriptチェック実行中...$(NC)"
	cd front && npx tsc --noEmit

# フロントエンドログを表示
frontend-logs:
	docker-compose logs -f frontend

# フロントエンドを再起動
restart-frontend:
	@echo "$(YELLOW)🔄 フロントエンドを再起動中...$(NC)"
	docker-compose restart frontend

# APIクライアント生成
api-generate:
	@echo "$(GREEN)🔄 APIクライアントを生成中...$(NC)"
	cd front && npm run api:generate

# ====================
# 🧹 クリーンアップ
# ====================

# ビルドキャッシュをクリア
clean:
	@echo "$(YELLOW)🧹 ビルドキャッシュをクリア中...$(NC)"
	# フロントエンドキャッシュクリア
	cd front && rm -rf build/ node_modules/.cache/ styled-system/
	# バックエンドキャッシュクリア
	cd backend && go clean -cache -modcache -testcache
	@echo "$(GREEN)✅ ビルドキャッシュをクリアしました$(NC)"

# 全キャッシュ・ボリューム・イメージを削除
clean-all: down
	@echo "$(RED)🗑️  全キャッシュ・ボリューム・イメージを削除中...$(NC)"
	# Node.jsキャッシュ削除
	cd front && rm -rf node_modules/ package-lock.json build/ styled-system/
	# Dockerイメージ削除
	docker-compose down --volumes --remove-orphans
	docker rmi $$(docker images | grep posting-app | awk '{print $$3}') 2>/dev/null || true
	# Dockerボリューム削除
	docker volume rm posting-app_postgres_data 2>/dev/null || true
	# Go キャッシュクリア
	cd backend && go clean -cache -modcache -testcache
	@echo "$(GREEN)✅ 全リソースを削除しました$(NC)"

# Node.js/Dockerキャッシュを完全削除
cache-clear:
	@echo "$(YELLOW)🧹 キャッシュをクリア中...$(NC)"
	# Node.jsキャッシュクリア
	cd front && npm cache clean --force 2>/dev/null || true
	cd front && rm -rf .eslintcache 2>/dev/null || true
	# Dockerビルドキャッシュクリア
	docker builder prune -f 2>/dev/null || true
	@echo "$(GREEN)✅ キャッシュをクリアしました$(NC)"

# 未使用Dockerリソースを削除
docker-clean:
	@echo "$(YELLOW)🐳 未使用Dockerリソースを削除中...$(NC)"
	docker system prune -f
	docker image prune -f
	docker container prune -f
	docker network prune -f

# Dockerボリュームを削除
volume-clean:
	@echo "$(YELLOW)📦 Dockerボリュームを削除中...$(NC)"
	docker volume prune -f

# 完全リセット
reset: clean-all frontend-install up
	@echo "$(GREEN)🎉 完全リセット完了！$(NC)"

# ====================
# ✅ 品質チェック
# ====================

# 全品質チェック実行
check-all: frontend-lint frontend-typecheck frontend-test backend-test
	@echo "$(GREEN)✅ 全品質チェック完了！$(NC)"

# ヘルスチェック実行
health:
	@echo "$(GREEN)🏥 ヘルスチェック実行中...$(NC)"
	@echo "📊 Docker containers:"
	@docker ps --filter "name=posting-app" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
	@echo ""
	@echo "🌐 Services status:"
	@curl -f http://localhost:3000/ >/dev/null 2>&1 && echo "✅ Frontend: OK" || echo "❌ Frontend: DOWN"
	@curl -f http://localhost:8080/health >/dev/null 2>&1 && echo "✅ Backend: OK" || echo "❌ Backend: DOWN"
	@docker exec posting-app-db pg_isready -U postgres >/dev/null 2>&1 && echo "✅ Database: OK" || echo "❌ Database: DOWN"

# ====================
# 📋 開発用ショートカット
# ====================

# 高速開発起動（キャッシュクリアなし）
dev-fast:
	@echo "$(GREEN)⚡ 高速開発起動中...$(NC)"
	docker-compose up

# ローカル開発（Dockerを使わない）
dev-local: db-up
	@echo "$(GREEN)💻 ローカル開発環境起動中...$(NC)"
	@echo "バックエンド: cd backend && go run main.go"
	@echo "フロントエンド: cd front && npm start"

# データベースのみリセット
db-reset:
	@echo "$(YELLOW)🗄️  データベースリセット中...$(NC)"
	docker-compose down
	docker volume rm posting-app_postgres_data 2>/dev/null || true
	$(MAKE) db-up db-migrate db-seed

# フロントエンドのみリセット
frontend-reset:
	@echo "$(YELLOW)🎨 フロントエンドリセット中...$(NC)"
	cd front && rm -rf node_modules/ package-lock.json build/ styled-system/
	$(MAKE) frontend-install
	docker-compose up -d --build frontend