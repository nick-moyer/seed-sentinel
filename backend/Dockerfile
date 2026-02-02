# ==========================================
# Stage 1: Base (Dependencies)
# ==========================================
FROM golang:1.25-alpine AS base
WORKDIR /app

# Install system dependencies required for SQLite (CGO)
RUN apk add --no-cache gcc musl-dev

# -----------------------------------------------
# [ADDED] Enable CGO globally for all stages
# This ensures 'go build' inside Air uses CGO for SQLite
ENV CGO_ENABLED=1
# -----------------------------------------------

# Cache Go Modules
COPY go.mod go.sum ./
RUN go mod download

# ==========================================
# Stage 2: Development (Hot Reload)
# ==========================================
FROM base AS dev
# Install Air for live reloading
RUN go install github.com/air-verse/air@latest

# -----------------------------------------------
# [ADDED] Copy config as a fallback
# Prevents crash if volume mount fails or lags
COPY .air.toml ./
# -----------------------------------------------

# Run Air (which runs your app)
CMD ["air", "-c", ".air.toml"]

# ==========================================
# Stage 3: Builder (Compile for Prod)
# ==========================================
FROM base AS builder
COPY . .
# CGO_ENABLED=1 is inherited from base, but specifying it here is fine too
RUN GOOS=linux go build -ldflags="-w -s" -o server .

# ==========================================
# Stage 4: Production (Minimal Runtime)
# ==========================================
FROM alpine:latest AS prod
WORKDIR /root/
# SQLite needs libc compatibility
RUN apk add --no-cache sqlite-libs
COPY --from=builder /app/server .
# Create data directory for volume mapping
RUN mkdir /data
ENV DB_PATH=/data/sentinel.db

EXPOSE 8080
CMD ["./server"]