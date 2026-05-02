# DailyTracker

[![GitHub release (latest by date)](https://img.shields.io/github/v/release/bsjhx/dailytracker-go?style=flat-square)](https://github.com/bsjhx/dailytracker-go/releases)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/bsjhx/dailytracker-go/build.yml?style=flat-square)](https://github.com/bsjhx/dailytracker-go/actions)
[![Go Version](https://img.shields.io/github/go-mod/go-version/bsjhx/dailytracker-go?style=flat-square)](go.mod)
[![License](https://img.shields.io/github/license/bsjhx/dailytracker-go?style=flat-square)](LICENSE)

A simple daily productivity tracker to rate your work and personal life (0-5 scale). Track your productivity, view statistics, and monitor your progress over time.

## ✨ Features

- 📊 Daily scoring for work and personal life (0-5)
- 📈 Weekly statistics and averages
- 🗂️ View last 30 entries
- ✏️ Edit existing entries
- 🚀 Fast, lightweight, single binary deployment
- 💾 PostgreSQL database
- 🐳 Docker support

## 🛠️ Tech Stack

- **Backend:** Go 1.25+ (pure Go, no CGO)
- **Database:** PostgreSQL
- **Frontend:** Vanilla JavaScript
- **Deployment:** Docker Compose

## 🚀 Quick Start

### Prerequisites

- Docker and Docker Compose (recommended)
- OR PostgreSQL 16+ installed locally

### Local Development with Docker

```bash
# Clone repository
git clone git@github.com:bsjhx/dailytracker-go.git
cd dailytracker-go

# Copy environment file
cp .env.example .env

# Start PostgreSQL and the app
docker-compose up

# Access app
open http://localhost:8080
```

### Local Development without Docker

```bash
# Make sure PostgreSQL is running
# Create database: dailytracker

# Set environment variables
export DB_URL=localhost:5432
export DB_USER=dailytracker
export DB_PASSWORD=your_password
export DB_NAME=dailytracker

# Run the application
go run ./cmd/dailytracker
```

### VPS Deployment

**First time setup:**
```bash
# On VPS
git clone git@github.com:bsjhx/dailytracker-go.git
cd dailytracker-go
./deploy-vps.sh
```

**Deploy new version (after code changes):**
```bash
# On VPS
cd dailytracker-go
./deploy-vps.sh
```

**Start/Stop existing deployment:**
```bash
# Start
./start.sh

# Stop
docker compose -f docker-compose.yml -f docker-compose.prod.yml down

# View logs
docker compose -f docker-compose.yml -f docker-compose.prod.yml logs -f
```

The app runs on **port 20224** on VPS, **port 8080** locally.

## 📖 Documentation

- **[DEPLOYMENT.md](DEPLOYMENT.md)** - Comprehensive deployment guide
- **[TODO.md](TODO.md)** - Project roadmap and tasks
- **[QUICKSTART.md](QUICKSTART.md)** - Quick reference guide

## 📡 API Endpoints

- `GET /api/entries` - Get last 30 entries
- `POST /api/entries` - Create new entry
- `GET /api/entries/:date` - Get entry by date (YYYY-MM-DD)
- `PUT /api/entries/:date` - Update entry
- `GET /api/stats/weekly` - Get 7-day statistics

## 📁 Project Structure

```
dailytracker-go/
├── api/                    # API handlers
│   ├── db.go              # Database connection
│   ├── entries.go         # Entries endpoints
│   ├── entry.go           # Single entry operations
│   └── stats.go           # Statistics endpoints
├── public/                # Frontend
│   └── index.html         # Single-page app
├── data/                  # Database storage (created on first run)
├── docker-compose.yml     # Local Docker config
├── docker-compose.prod.yml # VPS Docker config
├── Dockerfile             # Docker build instructions
├── deploy-vps.sh          # VPS deployment script
├── start.sh               # Start existing containers (VPS)
├── main.go                # Application entry point
└── go.mod                 # Go dependencies
```

## 🔧 Configuration

### Ports

- **Local:** `8080`
- **VPS:** `20224` (configurable in `docker-compose.prod.yml`)

### Database

PostgreSQL database with persistent storage:
- **Docker:** PostgreSQL data persisted in `postgres_data` volume
- **Port:** 5432 (accessible for direct connections)

### Environment Variables

Required for all environments:
```bash
DB_URL=localhost:5432       # PostgreSQL host:port
DB_USER=dailytracker        # Database user
DB_PASSWORD=your_password   # Database password
DB_NAME=dailytracker        # Database name
PORT=8080                   # App port (optional, default: 8080)
```

### Backup

```bash
# Backup database (Docker)
docker exec dailytracker-postgres pg_dump -U dailytracker dailytracker > backup-$(date +%Y%m%d).sql

# Restore (Docker)
docker exec -i dailytracker-postgres psql -U dailytracker dailytracker < backup-20260412.sql
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

MIT License - feel free to use this project however you'd like!

## 🙏 Acknowledgments

- PostgreSQL for reliable database
- Docker for easy deployment
- Go for being awesome
