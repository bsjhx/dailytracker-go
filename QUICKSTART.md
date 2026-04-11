# Quick Start Guide

## For VPS Deployment

### 1. Copy files to your VPS

```bash
scp -r dailytracker-go/ user@your-vps-ip:/home/user/
```

### 2. SSH into your VPS

```bash
ssh user@your-vps-ip
cd /home/user/dailytracker-go
```

### 3. Setup environment

```bash
cp .env.example .env
nano .env  # Set POSTGRES_PASSWORD
```

### 4. Deploy

```bash
./deploy.sh
```

That's it! Application will be available at `http://your-vps-ip:8080`

---

## For Local Testing

```bash
# 1. Set environment
cp .env.example .env
nano .env  # Set POSTGRES_PASSWORD

# 2. Run
./deploy.sh

# 3. Open browser
open http://localhost:8080
```

---

## Commands

```bash
# View logs
docker-compose logs -f

# Stop
docker-compose down

# Restart
docker-compose restart app

# Rebuild after code changes
docker-compose up -d --build
```

---

## For Production with Domain

See [DEPLOYMENT.md](DEPLOYMENT.md) for:
- Nginx reverse proxy setup
- HTTPS with Let's Encrypt
- Backup strategies
