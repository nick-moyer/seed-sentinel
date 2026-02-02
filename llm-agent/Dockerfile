# 1. Base Image
FROM python:3.12-slim-bookworm

# 2. Install uv
COPY --from=ghcr.io/astral-sh/uv:latest /uv /uvx /bin/

# 3. Setup Work Directory
WORKDIR /app

# 4. Copy Manifest Files
COPY pyproject.toml uv.lock ./

# 5. Install Dependencies
RUN uv sync --frozen --no-dev --no-install-project

# 6. Copy Application Code
COPY . .

# 7. Path Setup
ENV PATH="/app/.venv/bin:$PATH"

# 8. Run
CMD ["uv", "run", "agent.py"]
