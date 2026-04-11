# Daily Tracker - VPS Deployment Guide

This is a simple standalone Go application that uses SQLite as its database.

## Prerequisites

- Go 1.22 or higher installed on your VPS
- GCC compiler (required for SQLite CGO bindings)
- rsync installed on your local machine (for deployment script)

## Quick Deployment with rsync Script

The easiest way to deploy is using the provided `deploy.sh` script:

```bash
# Make the script executable (first time only)
chmod +x deploy.sh

# Deploy to your VPS
./deploy.sh user@your-vps-ip:/path/to/deployment

# Example:
./deploy.sh mikrus:/root/dailytracker 10224
```

The script will:
- Copy all necessary files to your VPS
- Exclude unnecessary files (.git, .db, Docker files, etc.)
- Show progress of the sync
- Provide next steps instructions

## Manual Installation on VPS

If you prefer to deploy manually:

1. Copy all files to your VPS location:
```bash
scp -r dailytracker-go user@your-vps:/path/to/deployment
```

Or using rsync:
```bash
rsync -avz --exclude '.git' --exclude '*.db' dailytracker-go/ user@your-vps:/path/to/deployment/
```

2. SSH into your VPS:
```bash
ssh user@your-vps
cd /path/to/deployment/dailytracker-go
```

3. Install dependencies (if needed):
```bash
go mod download
```

4. Run the application:
```bash
go run main.go
```

The server will start on port 8080 by default.

## Custom Port

To use a different port, set the PORT environment variable:
```bash
PORT=3000 go run main.go
```

## Running in Background

To run the application in the background:
```bash
nohup go run main.go > app.log 2>&1 &
```

Or compile and run:
```bash
go build -o dailytracker main.go
nohup ./dailytracker > app.log 2>&1 &
```

## Database

The application uses SQLite and stores data in `dailytracker.db` file in the same directory. This file is created automatically on first run.

## API Endpoints

- `GET /` - Serves static frontend from `public/` directory
- `GET /api/entries` - Get all daily entries (last 30)
- `POST /api/entries` - Create a new daily entry
- `GET /api/entries/:date` - Get entry by date
- `PUT /api/entries/:date` - Update entry by date
- `GET /api/stats/weekly` - Get weekly statistics

## Stopping the Application

Find the process ID:
```bash
ps aux | grep dailytracker
```

Kill the process:
```bash
kill <PID>
```
