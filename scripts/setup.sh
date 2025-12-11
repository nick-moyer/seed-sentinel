#!/bin/bash
set -e # Exit immediately if a command fails

echo "Seed Sentinel: Setup (Linux/Mac)"
echo "=========================================="

# 1. CHECK/INSTALL OLLAMA
if ! command -v ollama &> /dev/null; then
    echo "⬇Ollama not found. Installing..."
    curl -fsSL https://ollama.com/install.sh | sh
else
    echo "Ollama is installed."
fi

# 2. START OLLAMA & PULL MODEL
if ! pgrep -x "ollama" > /dev/null; then
    echo "Starting Ollama Service..."
    ollama serve &
    sleep 5
fi

if ! ollama list | grep -q "llama3"; then
    echo "Downloading Llama 3 Model (approx 4GB)..."
    ollama pull llama3
else
    echo "Llama 3 model is ready."
fi

# 3. CHECK/INSTALL UV
if ! command -v uv &> /dev/null; then
    echo "⬇uv not found. Installing..."
    curl -LsSf https://astral.sh/uv/install.sh | sh
    # Ensure uv is in path for this session
    source $HOME/.cargo/env 2>/dev/null || true
else
    echo "uv is installed."
fi

# 4. INSTALL PYTHON DEPS
echo "Installing Python Agent dependencies..."
cd llm-agent
uv sync
cd ..

# 5. INSTALL GO DEPS
echo "Installing Go Backend dependencies..."
cd backend
go mod tidy
cd ..

echo "=========================================="
echo "Setup Complete!"
echo "To run the system:"
echo "1. Terminal 1: cd llm-agent && uv run agent.py"
echo "2. Terminal 2: cd backend && go run main.go"