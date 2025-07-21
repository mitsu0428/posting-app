# 投稿プラットフォーム (Posting App)

サブスクリプション制の投稿プラットフォーム。Go製バックエンドとReact製フロントエンドで構築された、管理者による投稿承認機能付きのWebアプリケーションです。

## 📋 主要機能

### 🔐 認証・ユーザー管理
- ユーザー登録（メール認証）
- JWT認証（アクセストークン15分、リフレッシュトークン720時間）
- パスワードリセット機能
- 管理者・一般ユーザーのロール制御
- アカウント無効化・ユーザーBAN機能

### 📝 コンテンツ管理
- 投稿作成・編集・削除（承認ワークフロー付き）
- 投稿への返信（匿名・実名選択可）
- 画像アップロード機能（サムネイル生成）
- カテゴリシステム（色付きタグ、投稿あたり最大5個）
- いいね機能（リアルタイム集計）
- 論理削除システム（is_deletedフラグ）

### 👥 グループ機能
- グループ作成・編集・削除
- グループオーナー制度（権限管理）
- メンバー管理（表示ユーザー名による検索・追加）
- メンバー除名・自主退会機能
- グループ限定投稿（メンバーシップコンテンツ）
- 表示ユーザー名のユニーク制約（セキュリティ強化）

### 💳 サブスクリプション機能
- Stripe連携による定期課金
- コンテンツ作成にはアクティブなサブスクリプションが必要
- サブスクリプション状態のリアルタイム追跡
- Webhook対応とバッチ同期機能

### 🛡️ 管理機能
- 投稿承認・拒否ワークフロー
- ユーザー管理・BAN機能
- コンテンツモデレーションダッシュボード
- 管理者専用アクセス制御

## 🏗️ 技術スタック

### バックエンド
- **言語**: Go 1.21
- **フレームワーク**: Chi router
- **データベース**: PostgreSQL（SQLマイグレーション）
- **認証**: JWT + bcrypt
- **決済**: Stripe API
- **メール**: SendGrid
- **バリデーション**: go-playground/validator
- **ログ**: slog + zerolog（JSON形式）

### フロントエンド
- **フレームワーク**: React 18 + TypeScript
- **スタイリング**: PandaCSS + Material-UI（管理画面）
- **状態管理**: React Query + Context API
- **フォーム**: react-hook-form + zod validation
- **ファイルアップロード**: react-dropzone + compressorjs
- **決済**: Stripe React components
- **テスト**: React Testing Library

### インフラ・開発環境
- **API仕様**: OpenAPI 3.0 specification
- **コンテナ**: Docker + Docker Compose
- **データベース**: PostgreSQL 15
- **ファイルストレージ**: ローカルファイルシステム（設定可能）
- **デプロイ**: Cloud Run対応

## 🚀 クイックスタート

### 前提条件
- Docker & Docker Compose
- Node.js 18+（フロントエンドローカル開発）
- Go 1.21+（バックエンドローカル開発）

### 環境セットアップ

1. **リポジトリクローン**:
```bash
git clone <repository-url>
cd posting-app
```

2. **環境変数設定**:
```bash
cp .env.example .env
```

3. **`.env`ファイルの設定**:
```env
# データベース
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=posting_app

# JWT
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_ACCESS_DURATION=15m
JWT_REFRESH_DURATION=720h

# Stripe
STRIPE_API_KEY=sk_test_your_stripe_secret_key
STRIPE_PRICE_ID=price_your_stripe_price_id
STRIPE_WEBHOOK_SECRET=whsec_your_webhook_secret

# SendGrid
SENDGRID_API_KEY=your_sendgrid_api_key

# アプリケーション
BASE_URL=http://localhost:3000
PORT=8080
```

### Docker Compose での開発

```bash
# 全サービス起動
docker-compose up

# 個別サービス起動
docker-compose up postgres  # データベースのみ
docker-compose up backend   # バックエンドのみ
docker-compose up frontend  # フロントエンドのみ
```

**アクセスURL**:
- フロントエンド: http://localhost:3000
- バックエンドAPI: http://localhost:8080
- データベース: localhost:5432

### ローカル開発

#### バックエンド
```bash
cd backend
go mod tidy
go run main.go
```

#### フロントエンド
```bash
cd front
npm install
npm run generate-api  # OpenAPI仕様からAPIクライアント生成
npm start
```

## 👤 デフォルト管理者アカウント

**管理者ログイン情報**:
- **メール**: admin@example.com
- **パスワード**: admin123

**管理画面**: http://localhost:3000/admin

## 📁 プロジェクト構成

```
posting-app/
├── backend/                 # Go バックエンド
│   ├── domain/             # ドメインモデル（Post, User, Group, Category等）
│   ├── usecase/            # ビジネスロジック層
│   ├── repository/         # データアクセス層
│   ├── handler/            # HTTPハンドラー層
│   ├── infrastructure/     # 外部サービス実装（DB, JWT等）
│   ├── di/                 # 依存性注入コンテナ
│   ├── migrations/         # データベースマイグレーション
│   └── main.go            # エントリーポイント
│
├── front/                  # React フロントエンド
│   ├── src/
│   │   ├── pages/         # ページコンポーネント（Home, Groups, Admin等）
│   │   ├── components/    # 共通コンポーネント（Layout, Routes等）
│   │   ├── context/       # React Context（認証状態管理）
│   │   ├── generated/     # 自動生成APIクライアント
│   │   ├── utils/         # ユーティリティ関数
│   │   └── types/         # TypeScript型定義
│   ├── package.json
│   ├── tsconfig.json
│   └── orval.config.js    # API生成設定
│
├── api/
│   └── schema.yaml        # OpenAPI 3.0 仕様書
│
├── batch/                 # バッチ処理
│   └── subscription_batch.go
│
├── docker-compose.yml     # 開発環境設定
├── CLAUDE.md             # 開発者向けガイド
└── README.md             # このファイル
```

## 📊 データベース設計

### 主要テーブル
- **users** - ユーザーアカウント・サブスクリプション状態
- **posts** - ユーザー投稿・承認状態
- **replies** - 投稿への返信（匿名可）
- **categories** - カテゴリマスタ
- **post_categories** - 投稿-カテゴリ関連（多対多）
- **likes** - いいね機能
- **groups** - グループ情報
- **group_members** - グループメンバー関係
- **subscriptions** - Stripeサブスクリプション追跡
- **password_resets** - パスワードリセットトークン

### マイグレーション履歴
- `001_initial_schema.sql` - 基本テーブル作成
- `002_seed_data.sql` - 初期データ投入（管理者アカウント等）
- `003_add_new_features.sql` - 新機能追加（カテゴリ、いいね、グループ、論理削除）
- `004_add_display_name_unique.sql` - セキュリティ強化（表示ユーザー名ユニーク制約）

## 🔧 主要API エンドポイント

### 認証系
- `POST /auth/login` - ユーザー認証
- `POST /auth/register` - ユーザー登録
- `POST /auth/logout` - ログアウト
- `POST /auth/forgot-password` - パスワードリセット要求
- `POST /auth/reset-password` - パスワードリセット実行

### 投稿系
- `GET /posts` - 承認済み投稿一覧
- `POST /posts` - 新規投稿作成（サブスクリプション必須）
- `PUT /posts/{id}` - 投稿更新
- `DELETE /posts/{id}` - 投稿削除
- `POST /posts/{id}/replies` - 返信追加
- `POST /posts/{id}/like` - いいね切り替え

### グループ系
- `GET /groups` - ユーザーのグループ一覧
- `POST /groups` - グループ作成
- `PUT /groups/{id}` - グループ更新
- `DELETE /groups/{id}` - グループ削除
- `GET /groups/{id}/members` - メンバー一覧
- `POST /groups/{id}/members/by-name` - 表示ユーザー名でメンバー追加
- `DELETE /groups/{id}/members/{memberId}` - メンバー除名
- `POST /groups/{id}/leave` - グループ退会

### 管理者系
- `GET /admin/posts` - 投稿管理（承認待ち等）
- `POST /admin/posts/{id}/approve` - 投稿承認
- `POST /admin/posts/{id}/reject` - 投稿拒否
- `GET /admin/users` - ユーザー管理
- `POST /admin/users/{id}/ban` - ユーザーBAN

### その他
- `GET /categories` - カテゴリ一覧
- `GET /users/search` - ユーザー検索
- `GET /subscription/status` - サブスクリプション状態確認
- `POST /subscription/create-checkout-session` - Stripeチェックアウト作成

## 🧪 テスト・検証

### バックエンドテスト
```bash
cd backend
go test ./...
go test ./... -coverprofile=coverage.out
```

### フロントエンドテスト
```bash
cd front
npm test
npm run lint
npm run lint:fix
```

### API生成・更新
```bash
cd front
npm run generate-api  # OpenAPI仕様からTypeScriptクライアント生成
```

## 🚀 デプロイメント

### Cloud Run デプロイ

1. **プロダクションビルド**:
```bash
# バックエンド
cd backend
docker build -f Dockerfile.cloudrun -t backend-image .

# フロントエンド
cd front
npm run build
docker build -f Dockerfile.cloudrun -t frontend-image .
```

2. **環境変数設定**:
   - データベース接続情報
   - JWT秘密鍵（強力なランダム文字列）
   - Stripe APIキー・Webhook秘密鍵
   - SendGrid APIキー
   - BASE_URL（フロントエンドドメイン）

## ⚠️ 重要な実装課題と解決

### 1. Chi Router のルート優先順位
**問題**: `/{id}/members/by-name` と `/{id}/members` のパス競合
**解決**: より具体的なパスを先に定義

```go
// 修正前（404エラー発生）
r.Delete("/{id}", handlers.Post.DeleteGroup)
r.Get("/{id}/members", handlers.Post.GetGroupMembers)

// 修正後（正常動作）
r.Post("/{id}/members/by-name", handlers.Post.AddGroupMemberByDisplayName)
r.Delete("/{id}/members/{memberId}", handlers.Post.RemoveGroupMember)
r.Get("/{id}/members", handlers.Post.GetGroupMembers)
r.Delete("/{id}", handlers.Post.DeleteGroup)  // より一般的なパスは最後
```

### 2. ユーザー誤認識防止
**問題**: 表示ユーザー名の重複により、別ユーザーが追加される可能性
**解決**: display_name にUNIQUE制約追加

```sql
-- 既存重複データの自動解決
UPDATE users SET display_name = display_name || '_' || id::text 
WHERE id IN (
    SELECT id FROM (
        SELECT id, ROW_NUMBER() OVER (PARTITION BY display_name ORDER BY id) as rn 
        FROM users
    ) t WHERE rn > 1
);

-- UNIQUE制約追加
ALTER TABLE users ADD CONSTRAINT users_display_name_unique UNIQUE (display_name);
```

### 3. Docker コンテナキャッシュ問題
**問題**: バックエンドコード変更がコンテナに反映されない
**解決**: 強制リビルドコマンド

```bash
docker-compose up --build backend
# または
docker-compose down && docker-compose up
```

## 🔒 セキュリティ考慮事項

1. **環境変数**: 秘密情報をバージョン管理にコミットしない
2. **JWT秘密鍵**: 本番環境では強力なランダム文字列を使用
3. **HTTPS**: 本番環境では必ずHTTPS使用
4. **データベース**: 強力なパスワード設定・ネットワークアクセス制限
5. **ファイルアップロード**: ファイル形式・サイズの検証
6. **レート制限**: 適切なレート制限設定
7. **表示ユーザー名**: UNIQUE制約によるユーザー誤認識防止

## 🤝 開発・コントリビューション

### 新機能追加手順
1. **データベース設計**: マイグレーションファイル作成
2. **バックエンド実装**: ドメイン → リポジトリ → ユースケース → ハンドラー の順
3. **API仕様更新**: `api/schema.yaml` 更新
4. **フロントエンド**: API再生成 → UI実装
5. **テスト**: バックエンド・フロントエンドテスト実行

### よくある開発課題
- **Chi Router**: 具体的なパスを一般的なパスより先に定義
- **API生成**: OpenAPI仕様変更後は `npm run generate-api` 実行必須
- **Docker**: バックエンド変更時はコンテナ再起動・リビルド
- **TypeScript**: models ディレクトリからの直接インポート使用

## 📈 今後の拡張予定

### 機能拡張
- グループ権限管理強化（副管理者ロール）
- 投稿スケジューリング機能
- 高度な検索・フィルター機能
- 通知システム
- ファイル種別拡張（動画、PDF等）

### 技術改善
- パフォーマンス最適化（ページネーション、キャッシュ）
- セキュリティ強化（レート制限、CSRF対策）
- 監視・ログ強化
- 自動テスト拡充
- CI/CD パイプライン構築

## 📄 ライセンス

MIT License - 詳細は LICENSE ファイルを参照

## 📞 サポート

1. このREADMEドキュメントを確認
2. API仕様書（`api/schema.yaml`）を参照
3. 開発者向けガイド（`CLAUDE.md`）を確認
4. GitHub Issue で質問・報告
5. 開発チームへ連絡

---

このプラットフォームは Clean Architecture 原則に従い、レイヤー間の明確な分離とテスタビリティを重視して設計されています。継続的な機能拡張と保守性を考慮した構成となっています。