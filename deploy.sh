#!/bin/bash

set -e

# Configuration
PROJECT_ID=${PROJECT_ID:-"your-gcp-project-id"}
REGION=${REGION:-"us-central1"}
BACKEND_SERVICE="posting-app-backend"
FRONTEND_SERVICE="posting-app-frontend"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

echo_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

echo_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if gcloud is installed
if ! command -v gcloud &> /dev/null; then
    echo_error "gcloud CLI is not installed. Please install it first."
    exit 1
fi

# Check if PROJECT_ID is set
if [ "$PROJECT_ID" = "your-gcp-project-id" ]; then
    echo_error "Please set PROJECT_ID environment variable"
    exit 1
fi

echo_info "Starting deployment to GCP Project: $PROJECT_ID"

# Set the project
echo_info "Setting GCP project..."
gcloud config set project $PROJECT_ID

# Enable required APIs
echo_info "Enabling required APIs..."
gcloud services enable cloudbuild.googleapis.com
gcloud services enable run.googleapis.com
gcloud services enable containerregistry.googleapis.com
gcloud services enable sqladmin.googleapis.com

# Build and push backend image
echo_info "Building backend image..."
docker build -f backend/Dockerfile.cloudrun -t gcr.io/$PROJECT_ID/$BACKEND_SERVICE:latest ./backend

echo_info "Pushing backend image..."
docker push gcr.io/$PROJECT_ID/$BACKEND_SERVICE:latest

# Build and push frontend image
echo_info "Building frontend image..."
docker build -f front/Dockerfile.cloudrun -t gcr.io/$PROJECT_ID/$FRONTEND_SERVICE:latest ./front

echo_info "Pushing frontend image..."
docker push gcr.io/$PROJECT_ID/$FRONTEND_SERVICE:latest

# Deploy backend to Cloud Run
echo_info "Deploying backend to Cloud Run..."
gcloud run deploy $BACKEND_SERVICE \
    --image gcr.io/$PROJECT_ID/$BACKEND_SERVICE:latest \
    --region $REGION \
    --platform managed \
    --allow-unauthenticated \
    --memory 512Mi \
    --cpu 1 \
    --concurrency 100 \
    --max-instances 10 \
    --set-env-vars "PORT=8080"

# Get backend URL
BACKEND_URL=$(gcloud run services describe $BACKEND_SERVICE --region $REGION --format 'value(status.url)')
echo_info "Backend deployed at: $BACKEND_URL"

# Deploy frontend to Cloud Run
echo_info "Deploying frontend to Cloud Run..."
gcloud run deploy $FRONTEND_SERVICE \
    --image gcr.io/$PROJECT_ID/$FRONTEND_SERVICE:latest \
    --region $REGION \
    --platform managed \
    --allow-unauthenticated \
    --memory 256Mi \
    --cpu 1 \
    --concurrency 80 \
    --max-instances 5 \
    --port 8080

# Get frontend URL
FRONTEND_URL=$(gcloud run services describe $FRONTEND_SERVICE --region $REGION --format 'value(status.url)')
echo_info "Frontend deployed at: $FRONTEND_URL"

echo_info "Deployment completed successfully!"
echo_info "Frontend: $FRONTEND_URL"
echo_info "Backend: $BACKEND_URL"

echo_warn "Don't forget to:"
echo_warn "1. Set up Cloud SQL PostgreSQL instance"
echo_warn "2. Configure environment variables in Cloud Run services"
echo_warn "3. Set up Stripe webhooks pointing to: $BACKEND_URL/subscription/webhook"
echo_warn "4. Configure SendGrid API key for email functionality"