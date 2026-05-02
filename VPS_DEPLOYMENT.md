# Daily Tracker - VPS Deployment Guide

This is a Go application that uses PostgreSQL as its database.

## Prerequisites

- Docker and Docker Compose installed on your VPS
- Git installed on your VPS
- Port 20224 available (or your custom port)

## Quick Deployment with Docker Compose

The easiest way to deploy is using Docker Compose:

```bash
# 1. SSH into your VPS
ssh user@your-vps

# 2. Clone the repository
git clone git@github.com:bsjhx/dailytracker-go.git
cd dailytracker-go

# 3. Configure environment
cp .env.example .env
nano .env  # Set DB_PASSWORD and other values

# 4. Deploy using the deployment script
chmod +x deploy-vps.sh
./deploy-vps.sh
```

The script will:
- Pull the latest code
- Stop and remove old containers
- Build new Docker images
- Start PostgreSQL and the application
- Show the deployment status

## Manual Installation on VPS

If you prefer to deploy manually:

1. Clone or copy the repository to your VPS:
```bash
git clone git@github.com:bsjhx/dailytracker-go.git
cd dailytracker-go
```

2. Configure environment variables:
```bash
cp .env.example .env
nano .env
```

Set the following in `.env`:
```bash
DB_URL=postgres:5432
DB_USER=dailytracker
DB_PASSWORD=your_secure_password_here
DB_NAME=dailytracker
PORT=8080
```

3. Start the application:
```bash
docker-compose up -d
```

The server will start on port 8080 by default (or the port you configured).

## Custom Port

To use a different port, edit `docker-compose.prod.yml`:

```yaml
services:
  app:
    ports:
      - "YOUR_PORT:YOUR_PORT"
    environment:
      PORT: YOUR_PORT
```

Then deploy:
```bash
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

## Updating the Application

To deploy updates:

```bash
cd dailytracker-go
./deploy-vps.sh
```

Or manually:
```bash
git pull origin main
docker-compose down
docker-compose up -d --build
```

## Database

The application uses PostgreSQL with data stored in a Docker volume (`postgres_data`). This ensures data persists across container restarts.

### Backup Database

```bash
# Create backup
docker exec dailytracker-postgres pg_dump -U dailytracker dailytracker > backup-$(date +%Y%m%d).sql
```

### Restore Database

```bash
# Restore from backup
docker exec -i dailytracker-postgres psql -U dailytracker dailytracker < backup-20260501.sql
```

### Access Database

```bash
# Connect to PostgreSQL
docker exec -it dailytracker-postgres psql -U dailytracker -d dailytracker

# List tables
\dt

# View users
SELECT * FROM users;

# Exit
\q
```

## API Endpoints

- `GET /` - Serves static frontend
- `GET /api/entries` - Get all daily entries (last 30)
- `POST /api/entries` - Create a new daily entry
- `GET /api/entries/:date` - Get entry by date
- `PUT /api/entries/:date` - Update entry by date
- `GET /api/stats/weekly` - Get weekly statistics
- `POST /api/auth/login` - User login
- `POST /api/auth/logout` - User logout
- `POST /api/users/create` - Create new user

## Monitoring

### View Logs

```bash
# Application logs
docker-compose logs -f app

# PostgreSQL logs
docker-compose logs -f postgres

# Both
docker-compose logs -f
```

### Check Status

```bash
docker-compose ps
```

### Check Resource Usage

```bash
docker stats dailytracker-app dailytracker-postgres
```

## Stopping the Application

```bash
# Stop all containers
docker-compose down

# Stop and remove volumes (WARNING: deletes database)
docker-compose down -v
```

## Troubleshooting

### Container won't start

```bash
# Check logs
docker-compose logs

# Check if ports are available
sudo lsof -i :8080
sudo lsof -i :5432
```

### Database connection issues

```bash
# Verify PostgreSQL is running
docker-compose ps postgres

# Check database logs
docker-compose logs postgres

# Test connection
docker-compose exec postgres psql -U dailytracker -d dailytracker
```

### Reset everything

```bash
# Stop and remove all containers and volumes
docker-compose down -v

# Remove images
docker rmi dailytracker-go-app postgres:16-alpine

# Start fresh
docker-compose up -d
```

## Security Considerations

- Change the default `DB_PASSWORD` in production
- Use a reverse proxy (nginx) with HTTPS for production
- Configure firewall to only allow necessary ports
- Regular database backups
- Keep Docker and images updated

## Systemd Service (Optional)

For automatic startup on boot, create a systemd service:

```bash
sudo nano /etc/systemd/system/dailytracker.service
```

```ini
[Unit]
Description=DailyTracker Application
Requires=docker.service
After=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/path/to/dailytracker-go
ExecStart=/usr/bin/docker-compose up -d
ExecStop=/usr/bin/docker-compose down
TimeoutStartSec=0

[Install]
WantedBy=multi-user.target
```

Enable and start:
```bash
sudo systemctl enable dailytracker
sudo systemctl start dailytracker
```
