#!/bin/bash
set -e # Exit immediately if a command fails

echo "🌱 Seed Sentinel: Setup (Linux/Mac)"
echo "=========================================="

# 1. CHECK/INSTALL OLLAMA
if ! command -v ollama &> /dev/null; then
    echo "⬇️  Ollama not found. Installing..."
    curl -fsSL https://ollama.com/install.sh | sh
else
    echo "✅ Ollama is installed."
fi

# 2. START OLLAMA & PULL MODEL
if ! pgrep -x "ollama" > /dev/null; then
    echo "🔄 Starting Ollama Service..."
    ollama serve &
    sleep 5
fi

if ! ollama list | grep -q "llama3"; then
    echo "⬇️  Downloading Llama 3 Model (approx 4GB)..."
    ollama pull llama3
else
    echo "✅ Llama 3 model is ready."
fi

# 3. CHECK/INSTALL UV (Python Tooling)
if ! command -v uv &> /dev/null; then
    echo "⬇️  uv not found. Installing..."
    curl -LsSf https://astral.sh/uv/install.sh | sh
    source $HOME/.cargo/env 2>/dev/null || true
else
    echo "✅ uv is installed."
fi

# 4. CHECK NODE.JS (Required for Frontend)
if ! command -v npm &> /dev/null; then
    echo "❌ Node.js/npm not found!"
    echo "   Please install Node.js (v18+) manually from https://nodejs.org/"
    echo "   or use a version manager like nvm."
    exit 1
else
    echo "✅ Node.js is installed."
fi

# 5. INSTALL PYTHON DEPS
echo "📦 Installing Python Agent dependencies..."
cd llm-agent
uv sync
cd ..

# 6. INSTALL GO DEPS
echo "📦 Installing Go Backend dependencies..."
cd backend
go mod tidy
cd ..

# 7. BUILD FRONTEND
echo "📦 Building React Frontend..."
cd frontend
npm install      # Get dependencies
npm run build    # Compile to static HTML/JS in /dist
cd ..

echo "=========================================="
echo "🎉 Setup Complete!"
echo ""
echo "To run the system (use 3 separate terminals):"
echo "1. 🐍 Agent:    cd llm-agent && uv run agent.py"
echo "2. 🐹 Backend:  cd backend && go run ."
echo "3. 🚀 Frontend: cd frontend && npm dev"
echo ""
echo "🌍 Access the dashboard at: http://localhost:3000"