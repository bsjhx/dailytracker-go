#!/bin/bash

# Daily Tracker Startup Script

PORT=${PORT:-8080}

echo "Starting Daily Tracker on port $PORT..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go 1.22 or higher."
    exit 1
fi

# Build the application
echo "Building application..."
go build -o dailytracker main.go

if [ $? -ne 0 ]; then
    echo "Error: Build failed"
    exit 1
fi

# Run the application
echo "Starting server..."
PORT=$PORT ./dailytracker
