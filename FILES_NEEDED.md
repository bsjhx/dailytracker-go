# Files Required for VPS Deployment

## Essential Files (must copy)

### Application Code
- `main.go` - Main application entry point
- `api/db.go` - Database connection and migrations
- `api/entries.go` - Entries list and create endpoints
- `api/entry.go` - Single entry get/update endpoints
- `api/stats.go` - Weekly statistics endpoint

### Dependencies
- `go.mod` - Go module definition
- `go.sum` - Go module checksums

### Frontend
- `public/index.html` - Frontend UI (and any other files in public/)

### Scripts (optional but recommended)
- `start.sh` - Convenient startup script

### Documentation (optional)
- `VPS_DEPLOYMENT.md` - Deployment instructions
- `REFACTORING_SUMMARY.md` - Change summary

## Files NOT Needed (excluded by deploy.sh)

- `.git/` - Git repository data
- `.claude/` - Claude Code configuration
- `*.db` - Local database files (will be created on VPS)
- `.env` - Environment variables (not needed anymore)
- `Dockerfile` - Docker configuration
- `docker-compose.yml` - Docker Compose configuration
- `vercel.json` - Vercel configuration
- `.dockerignore` - Docker ignore file
- `DEPLOYMENT.md` - Old deployment docs
- `README.md` - Original readme
- `QUICKSTART.md` - Quickstart guide
- `.DS_Store` - macOS metadata
- `dailytracker` - Compiled binary (will be built on VPS)
- `*.exe`, `*.test` - Other binaries

## Total Size

Essential files are very small:
- Go source files: ~15-20 KB
- go.mod/go.sum: ~1 KB
- public/index.html: ~8 KB
- Scripts: ~2-3 KB

**Total: ~25-30 KB** (excluding dependencies which are downloaded via `go mod download`)

## What Gets Created on VPS

After deployment, the VPS will have:
- All source files
- `dailytracker.db` - SQLite database (created automatically)
- `dailytracker` - Compiled binary (if you use `go build`)
