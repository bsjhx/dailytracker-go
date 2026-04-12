# DailyTracker

A simple daily productivity tracker to rate your work and personal life (0-5 scale). Track your productivity, view statistics, and monitor your progress over time.

## ✨ Features

- 📊 Daily scoring for work and personal life (0-5)
- 📈 Weekly statistics and averages
- 🗂️ View last 30 entries
- ✏️ Edit existing entries
- 🚀 Fast, lightweight, single binary deployment
- 💾 SQLite database (no server required)
- 🐳 Docker support

## 🛠️ Tech Stack

- **Backend:** Go 1.25+ (pure Go, no CGO)
- **Database:** SQLite (via modernc.org/sqlite - pure Go implementation)
- **Frontend:** Vanilla JavaScript
- **Deployment:** Docker Compose

## 🚀 Quick Start

### Local Development

```bash
# Clone repository
git clone git@github.com:bsjhx/dailytracker-go.git
cd dailytracker-go

# Run directly with Go
go run main.go

# Or with Docker
docker compose up -d

# Access app
open http://localhost:8080
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

Database is stored as a single SQLite file:
- **Local:** `./dailytracker.db`
- **Docker:** `./data/dailytracker.db` (persisted volume)

### Backup

```bash
# Backup database
cp ./data/dailytracker.db backup-$(date +%Y%m%d).db

# Restore
cp backup-20260412.db ./data/dailytracker.db
docker compose restart
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

- Pure Go SQLite driver by [modernc.org/sqlite](https://gitlab.com/cznic/sqlite)
- Docker for easy deployment
- Go for being awesome
