#!/bin/bash

# Start DailyTracker on VPS (production) without rebuilding
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d

echo "✅ DailyTracker started in background on port 20224"
echo "📋 View logs: docker compose -f docker-compose.yml -f docker-compose.prod.yml logs -f"
echo "🛑 Stop: docker compose -f docker-compose.yml -f docker-compose.prod.yml down"
