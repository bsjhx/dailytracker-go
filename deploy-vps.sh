#!/bin/bash

set -e  # Exit on error

echo "🚀 Deploying DailyTracker..."

# 1. Update git repository
echo "📥 Pulling latest changes from git..."
git pull origin main

# 2. Stop and remove old container
echo "🛑 Stopping old container..."
docker-compose -f docker-compose.yml -f docker-compose.prod.yml down

# 3. Build new image
echo "🔨 Building new image..."
docker-compose -f docker-compose.yml -f docker-compose.prod.yml build --no-cache

# 4. Start new container in background
echo "✅ Starting new container..."
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

# 5. Show logs
echo "📋 Container started! Showing logs (Ctrl+C to exit)..."
sleep 2
docker-compose -f docker-compose.yml -f docker-compose.prod.yml logs -f
