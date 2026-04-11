# DailyTracker - Deployment Guide

Daily productivity tracker - rate your work and personal life (0-5).

## Stack
- **Backend:** Go with pure Go SQLite (no CGO required)
- **Database:** SQLite (single file, no server needed)
- **Frontend:** Vanilla JavaScript
- **Deployment:** Docker Compose

---

## Local Development

### Quick Start

```bash
# 1. Clone repository
git clone git@github.com:bsjhx/dailytracker-go.git
cd dailytracker-go

# 2. Run with Docker
docker-compose up -d

# 3. Access app
open http://localhost:8080
```

### Without Docker

```bash
# 1. Install Go 1.25+
# 2. Run
go run main.go

# App runs on http://localhost:8080
# Database created at ./dailytracker.db
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

# 4. Deploy
./deploy-vps.sh
```

The app will be available at `http://your-vps-ip:20224`

### Automated Deployment Script

The `deploy-vps.sh` script handles:
1. ✅ Pull latest code from git
2. ✅ Stop and remove old containers
3. ✅ Build new Docker image
4. ✅ Start new container in background

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
docker-compose -f docker-compose.yml -f docker-compose.prod.yml down

# Build and start
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build

# Check logs
docker-compose -f docker-compose.yml -f docker-compose.prod.yml logs -f
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

### Database

Database is stored in `./data/dailytracker.db` (persisted via Docker volume).

**Backup:**
```bash
# Copy database file
cp ./data/dailytracker.db backup-$(date +%Y%m%d).db
```

**Restore:**
```bash
# Replace database file
cp backup-20260411.db ./data/dailytracker.db
docker-compose restart
```

---

## Docker Commands

```bash
# View logs
docker-compose logs -f app

# Check status
docker-compose ps

# Restart app
docker-compose restart app

# Stop everything
docker-compose down

# Remove everything including database
docker-compose down -v

# Rebuild after code changes
docker-compose up -d --build

# Enter container shell
docker-compose exec app sh

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

# Change port in docker-compose.prod.yml
```

### Container won't start
```bash
# Check logs
docker-compose logs app

# Check if Docker is running
docker ps
```

### Database issues
```bash
# Check database file exists
ls -la ./data/

# Restart container
docker-compose restart app
```

### Can't connect from browser
```bash
# Check firewall
sudo ufw allow 20224

# Check container is running
docker-compose ps
```

---

## Files Structure

```
dailytracker-go/
├── api/              # API handlers
├── public/           # Frontend (index.html)
├── data/             # Database (created on first run)
├── docker-compose.yml           # Local config
├── docker-compose.prod.yml      # Production overrides
├── Dockerfile                   # Docker build
├── deploy-vps.sh               # Deployment script
├── main.go                     # Entry point
└── go.mod                      # Go dependencies
```

---

## Security Notes

- ✅ Database stored locally (not exposed)
- ✅ No environment secrets needed
- ✅ Pure Go SQLite (no C dependencies)
- ⚠️ No authentication (add if needed for public deployment)
- ⚠️ Consider HTTPS with nginx + Let's Encrypt for production
