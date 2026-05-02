.PHONY: help dev prod start stop restart deploy logs clean build test create-user db-only

# Default target
help:
	@echo "📋 DailyTracker - Available Commands"
	@echo ""
	@echo "Development:"
	@echo "  make dev              - Start in development mode (port 8080)"
	@echo "  make db-only          - Start only PostgreSQL database"
	@echo "  make build            - Build the application"
	@echo "  make test             - Run tests"
	@echo ""
	@echo "Production:"
	@echo "  make prod             - Start in production mode (port 20224)"
	@echo "  make deploy           - Full deployment (pull, build, start) - run on VPS"
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
dev:
	@echo "🚀 Starting DailyTracker in development mode..."
	docker compose up -d
	@echo "✅ DailyTracker started on http://localhost:8080"
	@echo "📋 View logs: make logs"

# Start only PostgreSQL
db-only:
	@echo "🗄️  Starting PostgreSQL only..."
	docker compose up postgres -d
	@echo "✅ PostgreSQL started on port 5432"

# Production mode (port 20224)
prod:
	@echo "🚀 Starting DailyTracker in production mode..."
	docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d
	@echo "✅ DailyTracker started on port 20224"
	@echo "📋 View logs: make logs-prod"
	@echo "🛑 Stop: make stop-prod"

# Start production (alias)
start: prod

# Stop development containers
stop:
	@echo "🛑 Stopping development containers..."
	docker compose down

# Stop production containers
stop-prod:
	@echo "🛑 Stopping production containers..."
	docker compose -f docker-compose.yml -f docker-compose.prod.yml down

# Restart development
restart: stop dev

# Restart production
restart-prod: stop-prod prod

# Full deployment (pull, build, start)
deploy:
	@echo "🚀 Deploying DailyTracker..."
	@echo "📥 Pulling latest changes from git..."
	git pull origin main
	@echo "🛑 Stopping old container..."
	docker compose -f docker-compose.yml -f docker-compose.prod.yml down
	@echo "🔨 Building new image..."
	docker compose -f docker-compose.yml -f docker-compose.prod.yml build --no-cache
	@echo "✅ Starting new container..."
	docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d
	@echo "📋 Container started! Showing logs (Ctrl+C to exit)..."
	@sleep 2
	docker compose -f docker-compose.yml -f docker-compose.prod.yml logs -f

# Show logs (development)
logs:
	docker compose logs -f

# Show logs (production)
logs-prod:
	docker compose -f docker-compose.yml -f docker-compose.prod.yml logs -f

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
