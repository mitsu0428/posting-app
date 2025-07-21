# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Backend (Go)
- `cd backend && go run main.go` - Run backend server
- `cd backend && go test ./...` - Run all tests
- `cd backend && go mod tidy` - Clean up dependencies

### Frontend (React + TypeScript)
- `cd front && npm install` - Install dependencies
- `cd front && npm start` - Start development server
- `cd front && npm run build` - Build for production
- `cd front && npm test` - Run tests
- `cd front && npm run lint` - Run linter
- `cd front && npm run lint:fix` - Fix linting errors
- `cd front && npm run generate-api` - Generate TypeScript API client from OpenAPI spec

### Full Stack Development
- `docker-compose up` - Start all services (database, backend, frontend)
- `docker-compose up postgres` - Start only PostgreSQL database
- `docker-compose down` - Stop all services

### API Code Generation
The frontend uses Orval to generate TypeScript API clients from the OpenAPI specification at `api/schema.yaml`. Always run `npm run generate-api` in the front directory after making changes to the API schema.

## Architecture Overview

This is a subscription-based posting platform with Clean Architecture principles:

### Backend Structure (Go)
- **Clean Architecture layers**:
  - `domain/` - Business entities (User, Post, Subscription)
  - `usecase/` - Business logic and application services
  - `repository/` - Data access layer interfaces and implementations
  - `handler/` - HTTP handlers and API endpoints
  - `infrastructure/` - External service implementations (database, JWT, email)
  - `di/` - Dependency injection container

- **Key dependencies**:
  - Chi router for HTTP routing
  - PostgreSQL with database/sql
  - JWT authentication with golang-jwt/jwt
  - Stripe SDK for payment processing
  - bcrypt for password hashing

### Frontend Structure (React + TypeScript)
- **Architecture**:
  - `src/pages/` - Route components and page layouts
  - `src/components/` - Reusable UI components
  - `src/context/` - React Context providers (AuthContext)
  - `src/generated/` - Auto-generated API client code (do not edit manually)
  - `src/utils/` - Utility functions and API configuration

- **Key technologies**:
  - React Query for server state management
  - Material-UI for admin interface components
  - PandaCSS for styling
  - react-hook-form + zod for form validation
  - Stripe React components for payment UI

### Database
PostgreSQL with migrations in `backend/migrations/`:
- `001_initial_schema.sql` - Core tables (users, posts, replies, subscriptions)
- `002_seed_data.sql` - Default admin user and test data
- `003_add_new_features.sql` - Advanced features (categories, likes, groups, logical deletion)
- `004_add_display_name_unique.sql` - Security enhancement (unique display names)

## Key Features and Business Logic

### Authentication & Authorization
- JWT-based auth with access (15m) and refresh (720h) tokens
- Role-based access control (admin vs regular users)
- Password reset flow with secure tokens

### Content Management
- Posts require subscription for creation
- Admin approval workflow for posts
- Anonymous and named replies
- Image upload with compression support
- Category system with color-coded tags (max 5 per post)
- Like/unlike functionality with real-time counts
- Logical deletion system (soft delete with is_deleted flag)

### Subscription System
- Stripe integration for recurring payments
- Webhook handling for real-time subscription updates
- Batch processing for webhook failure recovery

### Group Management
- Create and manage groups with names and descriptions
- Group ownership system with owner privileges
- Member management by display name (unique constraint for security)
- Member addition, removal (kick), and self-leave functionality
- Group-specific posts (membership-only content)
- Role-based permissions (owner vs member)

## Environment Configuration

Required environment variables (see docker-compose.yml for examples):
- `DB_*` - PostgreSQL connection details
- `JWT_SECRET` - Secret for signing JWT tokens
- `STRIPE_*` - Stripe API keys and webhook secrets
- `BASE_URL` - Frontend URL for redirects

## Common Development Patterns

### Backend
- All handlers follow dependency injection pattern via the DI container
- Repository pattern for data access with interface abstractions
- Use slog for structured logging
- Environment configuration via envconfig package
- Database transactions for multi-step operations

### Frontend
- All API calls use the generated client from `src/generated/api.ts`
- Global auth state managed via AuthContext
- React Query for caching and server state
- Form validation with react-hook-form + zod schemas

### Testing Strategy
- Backend: Unit tests for usecases and repositories
- Frontend: Component tests with React Testing Library
- Integration tests via Docker Compose setup

## Known Issues and Workarounds

### Backend Migration Issues
If you see "failed to run migrations: first .: file does not exist", the backend is trying to run database migrations but can't find the migration files. Ensure migrations are properly mounted in the Docker container.

### API Client Generation
After changing the OpenAPI schema (`api/schema.yaml`), you must regenerate the frontend API client by running `npm run generate-api` in the front directory before the frontend will work with the new API changes.

### Chi Router Route Precedence
Routes must be defined in order from most specific to most general. For example, `/{id}/members` must come before `/{id}` to avoid path conflicts. When adding new sub-routes, place them before general parameter routes.

### Docker Container Updates
When modifying backend code, the Docker container may cache old versions. Rebuild with `docker-compose up --build backend` or restart the container to ensure latest changes are picked up.

### Development Setup
For the full application to work properly:
1. Start PostgreSQL first: `docker-compose up postgres`
2. Run migrations by starting the backend: `docker-compose up backend`
3. Start frontend: `docker-compose up frontend` or `cd front && npm start`

## File Upload System
- Images are uploaded to `/app/uploads` in the backend container
- Frontend uses react-dropzone with compressorjs for client-side compression
- File validation is handled on both client and server side

---

# 詳細ファイル構成（日本語版）

## 1. バックエンド構成 (`backend/`)

### ドメイン層 (`domain/`)
- **`post.go`** - 投稿関連のエンティティ
  - `Post` - 投稿エンティティ（タイトル、コンテンツ、作者、ステータス、削除フラグ、グループID、カテゴリ、いいね数）
  - `Reply` - 返信エンティティ（コンテンツ、投稿ID、作者、匿名フラグ）
  - `Category` - カテゴリエンティティ（名前、説明、色）
  - `Like` - いいねエンティティ（投稿ID、ユーザーID）
  - `Group` - グループエンティティ（名前、説明、オーナー、メンバー）
  - `GroupMember` - グループメンバーエンティティ（グループID、ユーザーID、ロール）
  - `PostStatus` - 投稿ステータス定数（pending, approved, rejected）

- **`user.go`** - ユーザー関連のエンティティ
  - `User` - ユーザーエンティティ（メール、パスワード、表示名、ロール、サブスクリプション状態）
  - `UserRole` - ユーザーロール定数（user, admin）
  - `UserSubscriptionStatus` - サブスクリプション状態定数（active, inactive, past_due, canceled）
  - `PasswordReset` - パスワードリセットエンティティ

- **`subscription.go`** - サブスクリプション関連のエンティティ
  - `Subscription` - サブスクリプションエンティティ（StripeサブスクリプションID、状態、期間）

### ユースケース層 (`usecase/`)
- **`post_usecase.go`** - 投稿関連のビジネスロジック
  - 投稿作成（カテゴリ紐づけ、グループ制限チェック、サブスクリプション確認）
  - 投稿更新（カテゴリ更新、権限チェック）
  - 投稿削除（管理者または作者のみ、論理削除）
  - 投稿取得（承認済み、グループ制限、いいね状態含む）
  - 返信作成（サブスクリプション確認、匿名対応）
  - カテゴリ管理（取得、作成）
  - いいね機能（トグル式、重複防止）
  - グループ管理（作成、メンバー追加、投稿取得）

- **`auth_usecase.go`** - 認証関連のビジネスロジック
  - ユーザー登録（パスワードハッシュ化、デフォルトロール設定）
  - ユーザーログイン（認証情報確認、JWTトークン生成）
  - 管理者ログイン（管理者権限確認）
  - パスワードリセット（トークン生成、メール送信）

- **`subscription_usecase.go`** - サブスクリプション関連のビジネスロジック
  - Stripe連携（チェックアウトセッション作成、Webhook処理）
  - サブスクリプション状態管理

### リポジトリ層 (`repository/`)
- **`post_repository.go`** - 投稿関連のデータアクセス
  - 投稿CRUD操作（論理削除対応）
  - カテゴリCRUD操作
  - 投稿-カテゴリ関連操作（多対多関係）
  - いいね操作（ユニーク制約、トグル機能）
  - グループ操作（作成、メンバー管理、権限チェック）
  - 複雑な検索クエリ（承認済み投稿、グループ投稿、ユーザー投稿）
  - JOIN処理とN+1問題対策

- **`user_repository.go`** - ユーザー関連のデータアクセス
  - ユーザーCRUD操作
  - メール重複チェック
  - ロール別検索

- **`subscription_repository.go`** - サブスクリプション関連のデータアクセス
  - サブスクリプション状態管理
  - Stripe ID連携

- **`password_reset_repository.go`** - パスワードリセット関連のデータアクセス
  - リセットトークン管理
  - 有効期限チェック

### ハンドラー層 (`handler/`)
- **`post_handler.go`** - 投稿関連のHTTPハンドラー
  - 投稿CRUD API（マルチパートフォーム対応、ファイルアップロード）
  - カテゴリ管理API
  - いいね API（トグル機能）
  - グループ管理API
  - 権限チェック（管理者、作者、グループメンバー）
  - バリデーション（タイトル200文字、コンテンツ5000文字、カテゴリ5個まで）

- **`auth_handler.go`** - 認証関連のHTTPハンドラー
  - ユーザー登録、ログイン、ログアウト
  - パスワードリセット
  - プロフィール更新

- **`admin_handler.go`** - 管理者関連のHTTPハンドラー
  - 投稿承認・拒否
  - ユーザー管理（BAN機能）
  - 管理画面専用API

- **`subscription_handler.go`** - サブスクリプション関連のHTTPハンドラー
  - Stripeチェックアウトセッション作成
  - Webhook受信処理
  - サブスクリプション状態確認

- **`user_handler.go`** - ユーザー関連のHTTPハンドラー
  - プロフィール管理
  - パスワード変更
  - アカウント無効化

- **`middleware.go`** - ミドルウェア
  - JWT認証ミドルウェア（トークン検証、ユーザー情報抽出）
  - 管理者権限チェックミドルウェア
  - CORS設定
  - ロギング

- **`router.go`** - ルーティング設定
  - 全APIエンドポイントの定義
  - ミドルウェア適用
  - 権限別ルートグループ

- **`handler.go`** - 共通ハンドラーユーティリティ
  - レスポンス作成ヘルパー（JSON、エラー）
  - パラメータ抽出（URL、クエリ）
  - バリデーション共通処理

### インフラストラクチャ層 (`infrastructure/`)
- **`database.go`** - データベース接続管理
  - PostgreSQL接続設定
  - コネクションプール管理
  - マイグレーション実行

- **`jwt.go`** - JWT関連の実装
  - アクセストークン生成（15分有効）
  - リフレッシュトークン生成（720時間有効）
  - トークン検証
  - Claims構造体定義

### 依存性注入 (`di/`)
- **`container.go`** - DIコンテナ
  - 全依存関係の管理
  - インターフェース実装の注入
  - ライフサイクル管理

### マイグレーション (`migrations/`)
- **`001_initial_schema.sql`** - 初期スキーマ
  - users, posts, replies, subscriptions, password_resets テーブル
  - 基本インデックス設定

- **`002_seed_data.sql`** - シードデータ
  - 管理者ユーザー作成
  - テストデータ投入

- **`003_add_new_features.sql`** - 新機能追加
  - is_deleted カラム追加（論理削除）
  - categories, post_categories テーブル（カテゴリ機能）
  - likes テーブル（いいね機能）
  - groups, group_members テーブル（メンバーシップ機能）
  - 関連インデックス追加

- **`004_add_display_name_unique.sql`** - セキュリティ強化
  - display_name カラムにUNIQUE制約追加
  - ユーザー誤認防止のための重複排除処理
  - パフォーマンス向上のためのインデックス追加

### エントリーポイント
- **`main.go`** - アプリケーションエントリーポイント
  - 環境変数読み込み
  - データベース接続
  - DIコンテナ初期化
  - HTTPサーバー起動

## 2. フロントエンド構成 (`front/`)

### ページコンポーネント (`src/pages/`)
- **`Home.tsx`** - ホームページ
  - 投稿一覧表示（ページネーション対応）
  - カテゴリタグ表示
  - いいねボタン（ハート、数値表示）
  - 削除ボタン（管理者・作者のみ）
  - 認証状態別表示制御

- **`CreatePost.tsx`** - 投稿作成ページ
  - マルチパートフォーム（タイトル、コンテンツ、サムネイル）
  - カテゴリ選択UI（最大5個、色付きタグ）
  - グループ選択（メンバーシップ投稿）
  - ファイルアップロード（プレビュー機能、5MB制限、JPEG/PNG）
  - リアルタイムバリデーション（文字数制限表示）

- **`PostDetail.tsx`** - 投稿詳細ページ
  - 投稿本文表示
  - 返信表示・作成
  - いいね機能
  - 編集・削除ボタン（権限制御）

- **`Login.tsx`** - ログインページ
  - メール・パスワード認証
  - エラーハンドリング
  - 管理者ログインリンク

- **`Register.tsx`** - ユーザー登録ページ
  - ユーザー情報入力フォーム
  - バリデーション（メール形式、パスワード強度）

- **`MyPage.tsx`** - マイページ
  - ユーザー投稿一覧
  - プロフィール編集
  - サブスクリプション状態表示

- **`AdminDashboard.tsx`** - 管理者ダッシュボード
  - 投稿承認・拒否機能
  - ユーザー管理
  - 統計情報表示

- **`AdminLogin.tsx`** - 管理者ログインページ
  - 管理者専用認証

- **`Subscription.tsx`** - サブスクリプション管理ページ
  - 現在のサブスクリプション状態表示
  - Stripeチェックアウト連携
  - 認証チェック機能

- **`Groups.tsx`** - グループ管理ページ
  - グループ作成・編集・削除機能
  - メンバー管理（追加・除名・退会）
  - メンバー一覧表示（ロール別UI）
  - オーナー権限制御
  - リアルタイム状態管理

- **`ForgotPassword.tsx`** - パスワードリセット要求ページ
- **`ResetPassword.tsx`** - パスワードリセット実行ページ

### 共通コンポーネント (`src/components/`)
- **`Layout.tsx`** - 全体レイアウト
  - ヘッダー（ナビゲーション、ユーザーメニュー）
  - フッター
  - レスポンシブ対応

- **`PrivateRoute.tsx`** - 認証が必要なルートの保護
  - ログイン状態チェック
  - 未認証時のリダイレクト

- **`AdminRoute.tsx`** - 管理者権限が必要なルートの保護
  - 管理者権限チェック
  - 権限不足時のアクセス拒否

### コンテキスト (`src/context/`)
- **`AuthContext.tsx`** - 認証状態管理
  - ユーザー情報（ID、メール、ロール、サブスクリプション状態）
  - ログイン・ログアウト機能
  - JWTトークン管理（localStorage）
  - トークン有効期限チェック
  - 自動ログアウト機能

### 自動生成API (`src/generated/`)
- **`api.ts`** - APIクライアント関数
  - 全エンドポイントのTypeScript関数
  - React Query hooks（useQuery, useMutation）
  - 型安全なAPI呼び出し

- **`models/`** - TypeScript型定義
  - バックエンドエンティティと1:1対応
  - リクエスト・レスポンス型
  - OpenAPIスキーマから自動生成

### ユーティリティ (`src/utils/`)
- **`api.ts`** - API設定
  - Axiosインスタンス設定
  - ベースURL設定
  - エラーハンドリング

- **`api-mutator.ts`** - APIリクエスト変換
  - 認証ヘッダー自動付与
  - レスポンス変換
  - エラー処理

### 型定義 (`src/types/`)
- **`index.ts`** - アプリケーション固有の型定義
  - 生成された型の拡張
  - UI状態管理用の型

### 設定ファイル
- **`package.json`** - 依存関係とスクリプト
- **`tsconfig.json`** - TypeScript設定
- **`orval.config.js`** - API生成設定（OpenAPI → TypeScript）
- **`panda.config.ts`** - CSS-in-JS設定

## 3. API仕様 (`api/`)

### OpenAPIスキーマ
- **`schema.yaml`** - REST API仕様
  - 全エンドポイント定義（認証、投稿、管理者、サブスクリプション）
  - リクエスト・レスポンススキーマ
  - 認証方式（JWT Bearer）
  - バリデーションルール
  - エラーレスポンス定義

#### 主要エンドポイントグループ：
1. **認証系** (`/auth/*`)
   - 登録、ログイン、ログアウト、パスワードリセット

2. **投稿系** (`/posts/*`)
   - CRUD操作、いいね、返信作成

3. **カテゴリ系** (`/categories/*`)
   - 一覧取得、作成（管理者のみ）

4. **グループ系** (`/groups/*`)
   - グループ管理、メンバー管理、グループ投稿

5. **管理者系** (`/admin/*`)
   - 投稿承認、ユーザー管理

6. **サブスクリプション系** (`/subscription/*`)
   - Stripe連携、状態確認、Webhook

## 4. バッチ処理 (`batch/`)

- **`subscription_batch.go`** - サブスクリプション関連バッチ
  - Webhook失敗時の再処理
  - サブスクリプション状態の定期同期
  - 期限切れユーザーの状態更新

## 5. インフラ・設定ファイル

### Docker設定
- **`docker-compose.yml`** - 開発環境構成
  - PostgreSQL（データベース）
  - Backend（Goアプリケーション）
  - Frontend（React開発サーバー）
  - 環境変数設定

- **`Dockerfile`** (backend/) - バックエンドコンテナ
- **`Dockerfile`** (front/) - フロントエンドコンテナ

### その他
- **`README.md`** - プロジェクト概要と起動手順
- **`.gitignore`** - Git除外設定
- **環境変数テンプレート** - 必要な環境変数の一覧

## 開発フロー

### 新機能追加時の手順
1. **データベース設計**
   - マイグレーションファイル作成 (`backend/migrations/`)
   - エンティティ定義更新 (`backend/domain/`)

2. **バックエンド実装**
   - リポジトリ層実装 (`repository/`)
   - ユースケース層実装 (`usecase/`)
   - ハンドラー層実装 (`handler/`)
   - DIコンテナ更新 (`di/`)

3. **API仕様更新**
   - OpenAPIスキーマ更新 (`api/schema.yaml`)
   - フロントエンドAPI再生成 (`npm run generate-api`)

4. **フロントエンド実装**
   - ページコンポーネント作成/更新 (`src/pages/`)
   - 共通コンポーネント更新 (`src/components/`)
   - 型定義更新 (`src/types/`)

5. **テスト・検証**
   - バックエンドテスト (`go test ./...`)
   - フロントエンドLint (`npm run lint`)
   - 統合テスト

この構成により、各ファイルの役割と依存関係が明確になり、機能追加時の影響範囲を正確に把握できます。