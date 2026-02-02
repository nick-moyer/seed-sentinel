# ==========================================
# Seed Sentinel Control Plane
# ==========================================

# Variables
COMPOSE_BASE := docker-compose.yml
COMPOSE_DEV  := docker-compose.dev.yml

# Phony targets
.PHONY: help dev prod build-dev build-prod down logs clean check-ollama

# Default target
help:
	@echo "🌱 Seed Sentinel Makefile"
	@echo "--------------------------------"
	@echo "make dev          -> Check Ollama + Run Hot-Reload Mode"
	@echo "make prod         -> Check Ollama + Run Production Mode"
	@echo "make down         -> Stop all containers"

# ==========================================
# 🦙 Ollama Setup (Runs before dev/prod)
# ==========================================
check-ollama:
	@echo "🔍 Checking AI Service Status..."
	@# 1. Install if missing
	@if ! command -v ollama >/dev/null; then \
		echo "⬇️  Ollama not found. Installing..."; \
		curl -fsSL https://ollama.com/install.sh | sh; \
	else \
		echo "✅ Ollama binary found."; \
	fi

	@# 2. Start Service (Bind to 0.0.0.0 for Docker Access)
	@if ! pgrep -x "ollama" > /dev/null; then \
		echo "🔄 Starting Ollama Service (Background)..."; \
		OLLAMA_HOST=0.0.0.0 ollama serve > /dev/null 2>&1 & \
		echo "⏳ Waiting for startup..."; \
		sleep 5; \
	else \
		echo "✅ Ollama is running."; \
	fi

	@# 3. Pull Model
	@if ! ollama list | grep -q "llama3"; then \
		echo "⬇️  Downloading Llama 3 Model (This may take a while)..."; \
		ollama pull llama3; \
	else \
		echo "✅ Llama 3 model is ready."; \
	fi

# ==========================================
# Workflows
# ==========================================

# Note the dependency: 'dev' runs 'check-ollama' first
dev: check-ollama
	@echo "🚀 Starting Development Environment..."
	docker compose -f $(COMPOSE_BASE) -f $(COMPOSE_DEV) up --build

prod: check-ollama
	@echo "🏭 Starting Production Environment..."
	docker compose -f $(COMPOSE_BASE) up --build -d
	@echo "✅ Started. Tailing logs..."
	docker compose logs -f

# ==========================================
# Utilities
# ==========================================
down:
	@echo "🛑 Stopping containers..."
	docker compose down --remove-orphans

logs:
	docker compose logs -f

clean:
	docker system prune -f