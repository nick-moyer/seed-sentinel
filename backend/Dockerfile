# ==========================================
# Stage 1: Base (Dependencies)
# ==========================================
FROM golang:1.25-alpine AS base
WORKDIR /app

# Disable CGO for a portable, lightweight binary (not needed for lib/pq)
ENV CGO_ENABLED=0

# Cache Go Modules
COPY go.mod go.sum ./
RUN go mod download

# ==========================================
# Stage 2: Development (Hot Reload)
# ==========================================
FROM base AS dev
# Install Air for live reloading
RUN go install github.com/air-verse/air@latest

COPY .air.toml ./

# Run Air
CMD ["air", "-c", ".air.toml"]

# ==========================================
# Stage 3: Builder (Compile for Prod)
# ==========================================
FROM base AS builder
COPY . .
RUN GOOS=linux go build -ldflags="-w -s" -o server .

# ==========================================
# Stage 4: Production (Minimal Runtime)
# ==========================================
FROM alpine:latest AS prod
WORKDIR /root/
# Add CA certificates for secure connections (e.g., to cloud DBs)
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/server .

EXPOSE 8080
CMD ["./server"]
