# DailyTracker - Deployment Guide

Daily productivity tracker - rate your work and personal life (0-5).

## Stack
- **Backend:** Go
- **Database:** PostgreSQL
- **Frontend:** Vanilla JavaScript
- **Deployment:** Docker Compose

---

## Local Development

### Quick Start with Docker

```bash
# 1. Clone repository
git clone git@github.com:bsjhx/dailytracker-go.git
cd dailytracker-go

# 2. Copy environment file
cp .env.example .env

# 3. Run with Docker (starts PostgreSQL + app)
docker-compose up

# 4. Access app
open http://localhost:8080
```

### Without Docker

```bash
# 1. Install PostgreSQL 16+
# 2. Create database
psql -U postgres
CREATE DATABASE dailytracker;
CREATE USER dailytracker WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE dailytracker TO dailytracker;

# 3. Set environment variables
export DB_URL=localhost:5432
export DB_USER=dailytracker
export DB_PASSWORD=your_password
export DB_NAME=dailytracker

# 4. Run application
go run ./cmd/dailytracker

# App runs on http://localhost:8080
```

---

## VPS Deployment (Production)

### Prerequisites

1. **VPS with Docker installed**
2. **Git repository cloned on VPS**
3. **Port 20224 available** (or configure your own)

### Initial Setup

```bash
# 1. SSH into VPS
ssh your-vps -p YOUR_PORT

# 2. Install Docker (if not installed)
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER
# Log out and back in for group changes

# 3. Clone repository
git clone git@github.com:bsjhx/dailytracker-go.git
cd dailytracker-go

# 4. Configure environment
cp .env.example .env
nano .env  # Set DB_PASSWORD and other values

# 5. Deploy
./deploy-vps.sh
```

The app will be available at `http://your-vps-ip:20224`

### Automated Deployment Script

The `deploy-vps.sh` script handles:
1. ✅ Pull latest code from git
2. ✅ Stop and remove old containers
3. ✅ Build new Docker image
4. ✅ Start new containers in background

**Usage:**
```bash
cd /path/to/dailytracker-go
./deploy-vps.sh
```

### Manual Deployment Steps

If you prefer manual control:

```bash
# Pull latest changes
git pull origin main

# Stop old containers
docker-compose down

# Build and start (includes PostgreSQL)
docker-compose up -d --build

# Check logs
docker-compose logs -f
```

---

## Configuration

### Ports

- **Local:** Runs on port `8080`
- **VPS:** Runs on port `20224` (configured in `docker-compose.prod.yml`)

To change VPS port, edit `docker-compose.prod.yml`:
```yaml
ports:
  - "YOUR_PORT:YOUR_PORT"
environment:
  PORT: YOUR_PORT
```

### Environment Variables

Required for all environments:
```bash
DB_URL=localhost:5432       # PostgreSQL host:port
DB_USER=dailytracker        # Database user
DB_PASSWORD=your_password   # Database password
DB_NAME=dailytracker        # Database name
PORT=8080                   # App port (optional, default: 8080)
```

### Database

PostgreSQL database with persistent storage via Docker volume.

**Backup:**
```bash
# Create SQL dump
docker exec dailytracker-postgres pg_dump -U dailytracker dailytracker > backup-$(date +%Y%m%d).sql

# Or backup the entire data directory
docker run --rm --volumes-from dailytracker-postgres -v $(pwd):/backup alpine tar czf /backup/postgres-backup-$(date +%Y%m%d).tar.gz /var/lib/postgresql/data
```

**Restore:**
```bash
# From SQL dump
docker exec -i dailytracker-postgres psql -U dailytracker dailytracker < backup-20260411.sql

# Or restore data directory
docker-compose down
docker volume rm dailytracker-go_postgres_data
docker-compose up -d postgres
# Wait for postgres to start
docker run --rm --volumes-from dailytracker-postgres -v $(pwd):/backup alpine sh -c "cd / && tar xzf /backup/postgres-backup-20260411.tar.gz"
docker-compose restart postgres
docker-compose up -d
```

---

## Docker Commands

```bash
# View logs
docker-compose logs -f app
docker-compose logs -f postgres

# Check status
docker-compose ps

# Restart app
docker-compose restart app

# Stop everything
docker-compose down

# Remove everything including database volume
docker-compose down -v

# Rebuild after code changes
docker-compose up -d --build

# Enter container shell
docker-compose exec app sh
docker-compose exec postgres psql -U dailytracker -d dailytracker

# Check image size
docker images | grep dailytracker
```

---

## Nginx Reverse Proxy (Optional)

For HTTPS and domain setup:

```nginx
# /etc/nginx/sites-available/dailytracker
server {
    listen 80;
    server_name tracker.yourdomain.com;

    location / {
        proxy_pass http://localhost:20224;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Enable and get HTTPS:
```bash
sudo ln -s /etc/nginx/sites-available/dailytracker /etc/nginx/sites-enabled/
sudo certbot --nginx -d tracker.yourdomain.com
sudo systemctl reload nginx
```

---

## API Endpoints

- `GET /api/entries` - Get last 30 entries
- `POST /api/entries` - Create new entry
- `GET /api/entries/:date` - Get entry by date (YYYY-MM-DD)
- `PUT /api/entries/:date` - Update entry
- `GET /api/stats/weekly` - Get 7-day statistics

---

## Troubleshooting

### Port already in use
```bash
# Check what's using the port
sudo lsof -i :20224

# Change port in docker-compose.yml or docker-compose.prod.yml
```

### Container won't start
```bash
# Check logs
docker-compose logs app
docker-compose logs postgres

# Check if Docker is running
docker ps
```

### Database connection issues
```bash
# Check if PostgreSQL is running
docker-compose ps postgres

# Check PostgreSQL logs
docker-compose logs postgres

# Verify environment variables
docker-compose exec app env | grep DB_

# Test connection manually
docker-compose exec postgres psql -U dailytracker -d dailytracker
```

### Can't connect from browser
```bash
# Check firewall
sudo ufw allow 20224

# Check container is running
docker-compose ps
```

### Migration issues
```bash
# Check migration logs
docker-compose logs app | grep -i migration

# Connect to database and check schema
docker-compose exec postgres psql -U dailytracker -d dailytracker
\dt  # List tables
\d users  # Describe users table
```

---

## Files Structure

```
dailytracker-go/
├── cmd/dailytracker/    # Application entry point
├── internal/            # Internal packages
│   ├── handlers/        # HTTP handlers
│   ├── middleware/      # Authentication middleware
│   ├── models/          # Data models
│   └── repository/      # Database layer
├── migrations/          # Database migrations
├── web/                 # Frontend files
├── docker-compose.yml            # Docker config
├── docker-compose.prod.yml       # Production overrides
├── Dockerfile                    # Docker build
├── deploy-vps.sh                # Deployment script
└── go.mod                       # Go dependencies
```

---

## Security Notes

- ✅ Database stored in Docker volume (not exposed externally)
- ✅ Session-based authentication implemented
- ✅ PostgreSQL password required
- ⚠️ Set strong DB_PASSWORD in production
- ⚠️ Consider HTTPS with nginx + Let's Encrypt for production
- ⚠️ Review firewall rules to restrict database access
