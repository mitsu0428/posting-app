#!/bin/bash

# Google Cloud Run デプロイスクリプト
set -e

# カラー出力用
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# ログ関数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 環境変数チェック
check_env() {
    if [ -z "$1" ]; then
        log_error "Environment variable $2 is not set"
        exit 1
    fi
}

# 設定値
PROJECT_ID=${GOOGLE_CLOUD_PROJECT:-""}
REGION=${GOOGLE_CLOUD_REGION:-"us-central1"}
SERVICE_NAME_BACKEND="posting-app-backend"
SERVICE_NAME_FRONTEND="posting-app-frontend"

# 必要な環境変数をチェック
check_env "$PROJECT_ID" "GOOGLE_CLOUD_PROJECT"

log_info "🚀 Starting deployment to Google Cloud Run"
log_info "Project ID: $PROJECT_ID"
log_info "Region: $REGION"

# Google Cloud認証チェック
if ! gcloud auth list --filter=status:ACTIVE --format="value(account)" | grep -q .; then
    log_error "No active Google Cloud authentication found. Please run 'gcloud auth login'"
    exit 1
fi

# プロジェクト設定
log_info "📋 Setting Google Cloud project..."
gcloud config set project $PROJECT_ID

# Container RegistryのAPIを有効化
log_info "🔧 Enabling required APIs..."
gcloud services enable cloudbuild.googleapis.com
gcloud services enable run.googleapis.com
gcloud services enable containerregistry.googleapis.com

# バックエンドのビルドとデプロイ
log_info "🏗️  Building and deploying backend..."
cd backend
gcloud builds submit --tag gcr.io/$PROJECT_ID/$SERVICE_NAME_BACKEND .

log_info "🚀 Deploying backend to Cloud Run..."
gcloud run deploy $SERVICE_NAME_BACKEND \
    --image gcr.io/$PROJECT_ID/$SERVICE_NAME_BACKEND \
    --platform managed \
    --region $REGION \
    --allow-unauthenticated \
    --port 8080 \
    --memory 1Gi \
    --cpu 1 \
    --max-instances 10 \
    --set-env-vars "PORT=8080,GO_ENV=production"

# バックエンドのURLを取得
BACKEND_URL=$(gcloud run services describe $SERVICE_NAME_BACKEND --platform managed --region $REGION --format="value(status.url)")
log_info "✅ Backend deployed at: $BACKEND_URL"

cd ..

# フロントエンドのビルドとデプロイ
log_info "🏗️  Building and deploying frontend..."
cd front

# フロントエンド用の環境変数を設定してビルド
export REACT_APP_API_URL="$BACKEND_URL/api"

gcloud builds submit --tag gcr.io/$PROJECT_ID/$SERVICE_NAME_FRONTEND .

log_info "🚀 Deploying frontend to Cloud Run..."
gcloud run deploy $SERVICE_NAME_FRONTEND \
    --image gcr.io/$PROJECT_ID/$SERVICE_NAME_FRONTEND \
    --platform managed \
    --region $REGION \
    --allow-unauthenticated \
    --port 8080 \
    --memory 512Mi \
    --cpu 1 \
    --max-instances 5

# フロントエンドのURLを取得
FRONTEND_URL=$(gcloud run services describe $SERVICE_NAME_FRONTEND --platform managed --region $REGION --format="value(status.url)")
log_info "✅ Frontend deployed at: $FRONTEND_URL"

cd ..

log_info "🎉 Deployment completed successfully!"
log_info "📱 Frontend URL: $FRONTEND_URL"
log_info "🔧 Backend URL: $BACKEND_URL"
log_info ""
log_warn "⚠️  Don't forget to:"
log_warn "  1. Configure your database connection"
log_warn "  2. Set up Stripe webhooks pointing to: $BACKEND_URL/api/subscription/webhook"
log_warn "  3. Update CORS settings if needed"