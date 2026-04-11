#!/bin/bash

# Daily Tracker - VPS Deployment Script
# Usage: ./deploy.sh user@your-vps-ip:/path/to/deployment [port]

# Check if destination argument is provided
if [ -z "$1" ]; then
    echo "Usage: ./deploy.sh user@your-vps-ip:/path/to/deployment [port]"
    echo "Example: ./deploy.sh root@192.168.1.100:/opt/dailytracker"
    echo "Example with custom port: ./deploy.sh mikrus:/opt/dailytracker 20224"
    exit 1
fi

DESTINATION=$1
SSH_PORT=${2:-22}  # Default to port 22 if not specified

echo "═══════════════════════════════════════════════"
echo "  Daily Tracker - VPS Deployment"
echo "═══════════════════════════════════════════════"
echo "Destination: $DESTINATION"
echo "SSH Port: $SSH_PORT"
echo ""

# Check if rsync is installed
if ! command -v rsync &> /dev/null; then
    echo "Error: rsync is not installed. Please install rsync first."
    exit 1
fi

# Files and directories to sync
echo "Syncing files to VPS..."
rsync -avz --progress \
    -e "ssh -p $SSH_PORT" \
    --exclude '.git' \
    --exclude '.DS_Store' \
    --exclude '*.db' \
    --exclude '.env' \
    --exclude '.env.local' \
    --exclude 'dailytracker' \
    --exclude '*.exe' \
    --exclude '*.test' \
    --exclude '.vercel' \
    --exclude 'postgres_data' \
    --exclude '.claude' \
    --exclude 'node_modules' \
    --exclude 'Dockerfile' \
    --exclude 'docker-compose.yml' \
    --exclude 'vercel.json' \
    --exclude '.dockerignore' \
    --exclude 'DEPLOYMENT.md' \
    --exclude 'README.md' \
    --exclude 'QUICKSTART.md' \
    ./ "$DESTINATION/"

if [ $? -eq 0 ]; then
    echo ""
    echo "═══════════════════════════════════════════════"
    echo "  ✓ Deployment successful!"
    echo "═══════════════════════════════════════════════"
    echo ""
    echo "Next steps:"
    echo "1. SSH into your VPS:"
    if [ "$SSH_PORT" != "22" ]; then
        echo "   ssh -p $SSH_PORT ${DESTINATION%:*}"
    else
        echo "   ssh ${DESTINATION%:*}"
    fi
    echo ""
    echo "2. Navigate to deployment directory:"
    echo "   cd ${DESTINATION#*:}"
    echo ""
    echo "3. Install dependencies (first time only):"
    echo "   go mod download"
    echo ""
    echo "4. Run the application:"
    echo "   go run main.go"
    echo ""
    echo "   Or use the start script:"
    echo "   ./start.sh"
    echo ""
    echo "   Or build and run:"
    echo "   go build -o dailytracker main.go"
    echo "   ./dailytracker"
    echo ""
else
    echo ""
    echo "✗ Deployment failed!"
    exit 1
fi
