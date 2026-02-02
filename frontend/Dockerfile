# ==========================================
# Stage 1: Base (Dependencies)
# ==========================================
FROM node:20-alpine AS base
WORKDIR /app
COPY package*.json ./
RUN npm ci

# ==========================================
# Stage 2: Development (Vite Server)
# ==========================================
FROM base AS dev
# Expose the Vite port
EXPOSE 5173
# --host is required for Docker networking
CMD ["npm", "run", "dev", "--", "--host"]

# ==========================================
# Stage 3: Builder (Compile Assets)
# ==========================================
FROM base AS builder
COPY . .
# This generates the /dist folder
RUN npm run build

# ==========================================
# Stage 4: Production (Nginx)
# ==========================================
FROM nginx:alpine AS prod
# Copy build artifacts
COPY --from=builder /app/dist /usr/share/nginx/html
# Copy custom Nginx config (proxy logic)
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]