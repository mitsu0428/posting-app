# Google Cloud Run デプロイガイド

このガイドでは、掲示板アプリをGoogle Cloud Runにデプロイする方法を説明します。

## 📋 事前準備

### 1. Google Cloud Project の作成
```bash
# Google Cloud CLIがインストールされていることを確認
gcloud version

# 新しいプロジェクトを作成（オプション）
gcloud projects create your-project-id

# プロジェクトを設定
gcloud config set project your-project-id

# 認証
gcloud auth login
```

### 2. 必要なAPIの有効化
```bash
gcloud services enable cloudbuild.googleapis.com
gcloud services enable run.googleapis.com
gcloud services enable containerregistry.googleapis.com
gcloud services enable sqladmin.googleapis.com
```

### 3. Cloud SQL の設定
```bash
# Cloud SQL インスタンスを作成
gcloud sql instances create posting-app-db \
    --database-version=POSTGRES_15 \
    --tier=db-f1-micro \
    --region=us-central1

# データベースを作成
gcloud sql databases create posting_app --instance=posting-app-db

# ユーザーを作成
gcloud sql users create postgres --instance=posting-app-db --password=your-secure-password
```

## 🔧 環境変数の設定

### Secret Manager での機密情報管理
```bash
# 機密情報をSecret Managerに保存
echo -n "your-secure-db-password" | gcloud secrets create DB_PASSWORD --data-file=-
echo -n "your-jwt-secret-key" | gcloud secrets create JWT_SECRET --data-file=-
echo -n "sk_live_your_stripe_secret" | gcloud secrets create STRIPE_SECRET_KEY --data-file=-
echo -n "whsec_your_webhook_secret" | gcloud secrets create STRIPE_WEBHOOK_SECRET --data-file=-

# データベース接続情報
echo -n "/cloudsql/your-project-id:us-central1:posting-app-db" | gcloud secrets create DB_HOST --data-file=-
echo -n "postgres" | gcloud secrets create DB_USER --data-file=-
```

## 🚀 デプロイ方法

### 自動デプロイ（推奨）
```bash
# 環境変数を設定
export GOOGLE_CLOUD_PROJECT="your-project-id"
export GOOGLE_CLOUD_REGION="us-central1"

# デプロイスクリプトを実行
./deploy.sh
```

### 手動デプロイ

#### バックエンドのデプロイ
```bash
# バックエンドをビルド
cd backend
gcloud builds submit --tag gcr.io/your-project-id/posting-app-backend

# Cloud Runにデプロイ
gcloud run deploy posting-app-backend \
    --image gcr.io/your-project-id/posting-app-backend \
    --platform managed \
    --region us-central1 \
    --allow-unauthenticated \
    --port 8080 \
    --memory 1Gi \
    --cpu 1 \
    --max-instances 10 \
    --add-cloudsql-instances your-project-id:us-central1:posting-app-db \
    --set-env-vars "PORT=8080,GO_ENV=production" \
    --set-secrets "DB_HOST=DB_HOST:latest,DB_USER=DB_USER:latest,DB_PASSWORD=DB_PASSWORD:latest,JWT_SECRET=JWT_SECRET:latest,STRIPE_SECRET_KEY=STRIPE_SECRET_KEY:latest"
```

#### フロントエンドのデプロイ
```bash
# フロントエンドをビルド
cd front
gcloud builds submit --tag gcr.io/your-project-id/posting-app-frontend

# Cloud Runにデプロイ
gcloud run deploy posting-app-frontend \
    --image gcr.io/your-project-id/posting-app-frontend \
    --platform managed \
    --region us-central1 \
    --allow-unauthenticated \
    --port 8080 \
    --memory 512Mi \
    --cpu 1 \
    --max-instances 5
```

## 🔧 デプロイ後の設定

### 1. データベースマイグレーション
```bash
# Cloud Runインスタンスに接続してマイグレーションを実行
gcloud run services proxy posting-app-backend --port=8080 &
curl -X POST http://localhost:8080/admin/migrate
```

### 2. Stripe Webhook の設定
1. Stripe Dashboard にログイン
2. Webhooks セクションに移動
3. 新しいエンドポイントを追加: `https://your-backend-url/api/subscription/webhook`
4. 必要なイベントを選択:
   - `customer.subscription.created`
   - `customer.subscription.updated`
   - `customer.subscription.deleted`
   - `invoice.payment_succeeded`
   - `invoice.payment_failed`

### 3. CORS設定の確認
バックエンドのCORS設定で、フロントエンドのURLが許可されていることを確認してください。

## 📊 モニタリング

### ログの確認
```bash
# バックエンドのログ
gcloud logs read "resource.type=cloud_run_revision AND resource.labels.service_name=posting-app-backend" --limit=50

# フロントエンドのログ
gcloud logs read "resource.type=cloud_run_revision AND resource.labels.service_name=posting-app-frontend" --limit=50
```

### メトリクスの監視
Google Cloud Console の Cloud Run セクションで以下を監視:
- リクエスト数
- レスポンス時間
- エラー率
- メモリ使用量
- CPU使用率

## 🔒 セキュリティ考慮事項

1. **Secret Manager**: 機密情報は必ずSecret Managerに保存
2. **IAM**: 最小権限の原則に従ってサービスアカウントを設定
3. **VPC**: 可能であればVPCコネクタを使用してプライベートネットワーク内で通信
4. **SSL/TLS**: HTTPS通信が自動的に有効化されることを確認
5. **認証**: 必要に応じてCloud Identity and Access Managementを設定

## 💰 コスト最適化

1. **リソース制限**: メモリとCPUの上限を適切に設定
2. **オートスケーリング**: 最大インスタンス数を適切に設定
3. **リクエストベース課金**: 使用量に応じた課金モデルを活用
4. **Cloud SQL**: 必要に応じてより小さなインスタンスタイプを選択

## 🐛 トラブルシューティング

### よくある問題と解決方法

1. **データベース接続エラー**
   - Cloud SQLインスタンスが起動していることを確認
   - Cloud SQL Proxyの設定を確認
   - Secret Managerの権限を確認

2. **ビルドエラー**
   - Dockerfileの構文を確認
   - 依存関係のバージョンを確認
   - Cloud Build の権限を確認

3. **デプロイエラー**
   - Cloud Run APIが有効になっていることを確認
   - IAM権限を確認
   - リソース制限を確認

### サポートリソース
- [Google Cloud Run ドキュメント](https://cloud.google.com/run/docs)
- [Cloud SQL ドキュメント](https://cloud.google.com/sql/docs)
- [Secret Manager ドキュメント](https://cloud.google.com/secret-manager/docs)