.PHONY: all build backend frontend clean dev test docker-up docker-down install deps

PREFIX := build
BACKEND_OUT := $(PREFIX)/hudautomata

all: build

# Install dependencies
deps:
	@echo "📦 Installing Go dependencies..."
	go mod download
	go mod tidy
	@echo "📦 Installing Frontend dependencies..."
	cd frontend && bun install

# Build everything
build: backend frontend

# Build backend
backend:
	@echo "🔨 Building backend..."
	mkdir -p "$(PREFIX)"
	CGO_ENABLED=1 go build -ldflags="-s -w" -o "$(BACKEND_OUT)" src/main.go
	@echo "✅ Backend built successfully: $(BACKEND_OUT)"

# Build frontend
frontend:
	@echo "🎨 Building frontend..."
	cd frontend && bun run build
	@echo "✅ Frontend built successfully"

# Development mode
dev:
	@echo "🚀 Starting development servers..."
	@make -j2 dev-backend dev-frontend

dev-backend:
	@echo "🔧 Starting backend (development)..."
	go run src/main.go

dev-frontend:
	@echo "🎨 Starting frontend (development)..."
	cd frontend && bun run dev

# Run tests
test:
	@echo "🧪 Running tests..."
	go test -v ./...

# Docker commands
docker-up:
	@echo "🐳 Starting Docker containers..."
	docker-compose up -d

docker-down:
	@echo "🐳 Stopping Docker containers..."
	docker-compose down

docker-build:
	@echo "🐳 Building Docker images..."
	docker-compose build

docker-logs:
	docker-compose logs -f

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf $(PREFIX)
	rm -rf frontend/dist
	rm -f hudautomata.db
	@echo "✅ Clean complete"

# Database operations
db-migrate:
	@echo "📊 Running database migrations..."
	go run src/main.go migrate

db-seed:
	@echo "🌱 Seeding database..."
	go run src/main.go seed

# Help
help:
	@echo "Available commands:"
	@echo "  make deps          - Install all dependencies"
	@echo "  make build         - Build backend and frontend"
	@echo "  make backend       - Build backend only"
	@echo "  make frontend      - Build frontend only"
	@echo "  make dev           - Start development servers"
	@echo "  make test          - Run tests"
	@echo "  make docker-up     - Start Docker containers"
	@echo "  make docker-down   - Stop Docker containers"
	@echo "  make clean         - Clean build artifacts"

