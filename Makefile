.PHONY: help dev prod start stop restart logs clean build test create-user db-only start-prod stop-prod

# Default target
help:
	@echo "📋 DailyTracker - Available Commands"
	@echo ""
	@echo "Development:"
	@echo "  make dev-docker       - Start in development mode (port 8080) - using docker"
	@echo "  make dev-bare         - Start in development mode (port 8080) - using go run"
	@echo "  make db-only          - Start only PostgreSQL database"
	@echo "  make build            - Build the application"
	@echo "  make test             - Run tests"
	@echo ""
	@echo "Production:"
	@echo "  make start-prod       - Start production (git pull, build, run on VPS)"
	@echo "  make stop-prod        - Stop production container"
	@echo "  make logs-prod        - Show production logs"
	@echo ""
	@echo "Docker Management:"
	@echo "  make stop             - Stop all containers"
	@echo "  make restart          - Restart containers"
	@echo "  make logs             - Show container logs"
	@echo "  make clean            - Stop and remove containers, volumes"
	@echo ""
	@echo "User Management:"
	@echo "  make create-user USER=<username> PASS=<password> [HOST=<host>]"
	@echo "    - Create a new user"
	@echo "    - Example: make create-user USER=admin PASS=secret"
	@echo "    - Example: make create-user USER=admin PASS=secret HOST=example.com:20224"

# Development mode (port 8080)
dev-docker:
	@echo "🚀 Starting DailyTracker in development mode..."
	docker compose up -d
	@echo "✅ DailyTracker started on http://localhost:8080"
	@echo "📋 View logs: make logs"

dev-bare:
	@echo "🚀 Starting DailyTracker in development mode (bare)..."
	go run cmd/dailytracker/main.go

# Start only PostgreSQL
db-only:
	@echo "🗄️  Starting PostgreSQL only..."
	docker compose up postgres -d
	@echo "✅ PostgreSQL started on port 5432"

# Stop development containers
down:
	@echo "🛑 Stopping development containers..."
	docker compose down

# Stop production container
stop-prod:
	@echo "🛑 Stopping production container..."
	@docker stop dailytracker-app 2>/dev/null || true
	@docker rm dailytracker-app 2>/dev/null || true
	@echo "✅ Production container stopped"

# Restart development
restart: stop dev

# Start production (git pull, build, restart with external Postgres)
start-prod:
	@echo "🚀 Starting Production (3 steps)..."
	@echo ""
	@echo "📥 Step 1/3: Pulling latest changes..."
	git pull origin main
	@echo ""
	@echo "🔨 Step 2/3: Building Docker image..."
	docker build -t dailytracker:latest .
	@echo ""
	@echo "🔄 Step 3/3: Starting container..."
	@docker stop dailytracker-app 2>/dev/null || true
	@docker rm dailytracker-app 2>/dev/null || true
	@docker run -d \
		--name dailytracker-app \
		--restart unless-stopped \
		--env-file .env \
		-p 20224:20224 \
		dailytracker:latest
	@echo ""
	@echo "✅ Production started!"
	@echo "📋 View logs: docker logs -f dailytracker-app"
	@echo "📊 Status: docker ps --filter name=dailytracker-app"

# Show logs (development)
logs:
	docker compose logs -f

# Show logs (production)
logs-prod:
	@docker logs -f dailytracker-app

# Clean up containers and volumes
clean:
	@echo "🧹 Cleaning up containers and volumes..."
	docker compose down -v
	@echo "✅ Cleanup complete"

# Build the application
build:
	@echo "🔨 Building DailyTracker..."
	go build -o bin/dailytracker ./cmd/dailytracker
	@echo "✅ Build complete: bin/dailytracker"

# Run tests
test:
	@echo "🧪 Running tests..."
	go test -v ./...

# Create a new user
create-user:
	@if [ -z "$(USER)" ] || [ -z "$(PASS)" ]; then \
		echo "❌ Error: USER and PASS are required"; \
		echo ""; \
		echo "Usage: make create-user USER=<username> PASS=<password> [HOST=<host>]"; \
		echo ""; \
		echo "Example:"; \
		echo "  make create-user USER=admin PASS=mypassword"; \
		echo "  make create-user USER=admin PASS=mypassword HOST=example.com:8080"; \
		exit 1; \
	fi
	@HOST=$${HOST:-localhost:8080}; \
	echo "Creating user '$(USER)' on $$HOST..."; \
	echo ""; \
	response=$$(curl -s -w "\n%{http_code}" -X POST "http://$$HOST/api/users/create" \
		-H "Content-Type: application/json" \
		-d "{\"username\":\"$(USER)\",\"password\":\"$(PASS)\"}"); \
	http_code=$$(echo "$$response" | tail -n1); \
	body=$$(echo "$$response" | sed '$$d'); \
	if [ "$$http_code" -eq 200 ] || [ "$$http_code" -eq 201 ]; then \
		echo "✅ User created successfully!"; \
		echo ""; \
		echo "You can now login with:"; \
		echo "  Username: $(USER)"; \
		echo "  Password: $(PASS)"; \
	else \
		echo "❌ Error creating user (HTTP $$http_code)"; \
		echo ""; \
		echo "Response:"; \
		echo "$$body"; \
		exit 1; \
	fi
