# 掲示板アプリ - 完全会員制掲示板システム

完全会員制の掲示板アプリです。ユーザー認証、投稿管理、サブスクリプション決済、管理者機能を備えた本格的なウェブアプリケーションです。

## 📋 主な機能

### 📝 会員登録機能
- ✅ 新規登録
- ✅ ログイン（`/login`）
- ✅ ログアウト
- ✅ パスワードを忘れた時の動線（`/forgot-password`、`/reset-password`）
- ✅ 管理者ログイン（`/admin-login-page`）

### 投稿機能
- **スレッド作成**
  - タイトル、サムネイル画像、内容での投稿
  - 管理者による承認制
- **返信機能**
  - 匿名投稿対応
  - ユーザー名表示選択可能

### マイページ機能
- 自分の投稿一覧表示
- 投稿ステータス確認（承認待ち・承認済み・却下）

### 管理者機能
- スレッド承認・却下・削除
- ユーザー一覧表示
- ユーザーアカウント無効化

### 決済機能
- Stripe決済によるサブスクリプション
- サブスクリプション状態の自動管理
- バッチ処理による定期的なステータス確認

## 技術スタック

### フロントエンド
- React 18
- TypeScript
- React Router Dom
- Axios
- Stripe.js

### バックエンド
- Go 1.21
- Gorilla Mux
- PostgreSQL
- JWT認証
- Stripe API
- Clean Architecture（Handler → Usecase → Repository）
- Dependency Injection（uber-go/dig）

### インフラ
- Docker & Docker Compose
- PostgreSQL 15
- Nginx

## セットアップ

### 前提条件
- Docker & Docker Compose
- Node.js 18+
- Go 1.21+

### 環境変数設定

`.env`ファイルを作成し、以下の環境変数を設定してください：

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=posting_app

# JWT
JWT_SECRET=your-jwt-secret-key

# Stripe
STRIPE_SECRET_KEY=sk_test_your_stripe_secret_key
STRIPE_WEBHOOK_SECRET=whsec_your_webhook_secret
STRIPE_PRICE_ID=price_your_price_id

# Frontend
REACT_APP_API_URL=http://localhost:8080/api
REACT_APP_STRIPE_PUBLISHABLE_KEY=pk_test_your_stripe_publishable_key
FRONTEND_URL=http://localhost:3000
```

### 起動方法

#### 🚀 簡単起動（推奨）- Makefileを使用
```bash
# ヘルプを表示
make help

# 開発環境を起動（全サービス）
make dev

# バックグラウンドで起動
make up

# 停止
make down

# キャッシュを完全クリアして起動
make reset
```

#### 📋 主なMakefileコマンド
```bash
# 開発環境管理
make dev              # 開発環境起動（フォアグラウンド）
make up               # 全サービス起動（バックグラウンド）
make down             # 全サービス停止
make restart          # 全サービス再起動
make status           # サービス状態確認

# データベース管理
make db-up            # DBのみ起動
make db-migrate       # マイグレーション実行
make db-seed          # シードデータ投入
make db-shell         # DB接続

# キャッシュクリア
make clean            # ビルドキャッシュクリア
make clean-all        # 全キャッシュ・イメージ削除
make cache-clear      # Node.js/Dockerキャッシュクリア
make reset            # 完全リセット

# 品質チェック
make check-all        # 全チェック実行
make health           # ヘルスチェック
```

#### 📋 手動起動（従来の方法）
```bash
# 1. 環境変数設定
cp .env.example .env
# 必要に応じて .env を編集（Stripeキーなど）

# 2. アプリケーション起動
./start-dev.sh
```

#### 📋 手動起動

1. **環境変数設定**
```bash
cp .env.example .env
# .envファイルを編集してStripeキーなどを設定
```

2. **データベース起動**
```bash
docker-compose up -d postgres
```

3. **バックエンド起動**
```bash
cd backend
go mod tidy
go run main.go
```

4. **フロントエンド（開発用）**
```bash
cd front
npm install --legacy-peer-deps
npm start
```

### アクセス
- フロントエンド: http://localhost:3000
- バックエンドAPI: http://localhost:8080/api
- 管理者ログイン: http://localhost:3000/admin-login-page

### デフォルト管理者アカウント
- Email: admin@example.com
- Password: admin123

## API仕様

OpenAPIスキーマは`/api/schema.yaml`に定義されています。

### 主要エンドポイント

#### 認証
- `POST /api/auth/register` - ユーザー登録
- `POST /api/auth/login` - ログイン
- `POST /api/admin/login` - 管理者ログイン

#### 投稿
- `GET /api/posts` - 投稿一覧取得
- `POST /api/posts` - 投稿作成
- `GET /api/posts/{id}` - 投稿詳細取得
- `GET /api/posts/{id}/replies` - 返信一覧取得
- `POST /api/posts/{id}/replies` - 返信作成

#### 管理者
- `GET /api/admin/posts` - 全投稿管理
- `POST /api/admin/posts/{id}/approve` - 投稿承認
- `POST /api/admin/posts/{id}/reject` - 投稿却下
- `DELETE /api/admin/posts/{id}` - 投稿削除

#### サブスクリプション
- `POST /api/subscription/create-checkout-session` - 決済セッション作成
- `POST /api/subscription/webhook` - Stripeウェブフック

## バッチ処理

サブスクリプション状態の確認バッチは以下で実行：

```bash
cd backend/batch
go run subscription_batch.go
```

## 開発

### テスト実行

バックエンド：
```bash
cd backend
go test ./...
```

フロントエンド：
```bash
cd front
npm test
```

### 型チェック・Lint

```bash
cd front
npm run typecheck
npm run lint
```

## デプロイ

本番環境では以下の設定を推奨：

1. 環境変数の適切な設定
2. HTTPS の有効化
3. データベースの適切な設定
4. Stripe本番キーの使用
5. セキュリティヘッダーの設定

## ライセンス

MIT License