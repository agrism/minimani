#!/usr/bin/env bash
set -e

echo "🚀 Starting MiniMani deployment..."

# 1. Pull latest git changes
git pull origin main

# 2. Build backend
echo "📦 Building backend..."
cd backend
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/server-linux-amd64 cmd/server/main.go
cd ..

# 3. Restart Docker containers
echo "🐳 Updating docker containers..."
docker compose up -d --build

# 4. Build frontend
echo "🎨 Building frontend..."
cd frontend
npm ci --prefer-offline || npm install
npm run build
cd ..

echo "✅ MiniMani deployed successfully to https://minimani.kilograms.lv"

