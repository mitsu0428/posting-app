# Node.js の公式イメージ（セキュリティアップデート版）
FROM node:20-alpine AS build

# 作業ディレクトリ作成
WORKDIR /app

# package.jsonとpackage-lock.jsonをコピー
COPY package*.json ./

# 依存関係をインストール
RUN npm ci --legacy-peer-deps

# アプリの全ファイルをコピー（Pandaの設定ファイルを含む）
COPY . .

# OpenAPI schema をコピー
COPY ../api ./api

# API生成とビルド
RUN npm run api:generate
RUN npm run panda
RUN npm run build

# ============================
# 本番サーバー（Nginx）で配信
# ============================
FROM nginx:alpine

# React Router のルーティング対応: fallback を index.html に
COPY ../front/nginx.conf /etc/nginx/nginx.conf

# ビルド成果物を Nginx に配置
COPY --from=build /app/build /usr/share/nginx/html

EXPOSE 3000

CMD ["nginx", "-g", "daemon off;"]
