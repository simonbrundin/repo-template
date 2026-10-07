# syntax=docker/dockerfile:1
#
# Nuxt SPA + Static File Server
#
# Architecture:
#   Browser → Nuxt SPA → Go API → PostgreSQL
#
# Stages:
#   - dev: Development server (bun)
#   - builder: Builds static SPA files
#   - runtime: nginx serves static files
#
# Usage:
#   Development: docker build --target dev -t nuxt-dev . && docker run -it -p 3000:3000 nuxt-dev
#   Production:  docker build -t nuxt-prod . && docker run -p 3000:80 nuxt-prod
#

# =============================================================================
# Dev Stage: Full development environment
# =============================================================================
FROM oven/bun:1-alpine AS dev

WORKDIR /app

# Install dependencies
COPY package.json bun.lock ./
RUN bun install --frozen-lockfile

# Copy source
COPY . .

EXPOSE 3000

CMD ["bun", "run", "dev", "--", "--host", "0.0.0.0"]

# =============================================================================
# Builder Stage: Build static SPA files
# =============================================================================
FROM oven/bun:1-alpine AS builder

WORKDIR /app

# Install dependencies
COPY package.json bun.lock ./
RUN bun install --frozen-lockfile

# Copy source and build
COPY . .
RUN bun run generate  # Generates static files in .output/public

# =============================================================================
# Runtime Stage: Serve static files with nginx
# =============================================================================
FROM nginx:alpine AS production

# Copy generated static files
COPY --from=builder /app/.output/public /usr/share/nginx/html

# Copy nginx config for SPA fallback
COPY --chown=nginx:nginx <<'EOF' /etc/nginx/conf.d/default.conf
server {
    listen 80;
    server_name _;
    root /usr/share/nginx/html;
    index index.html;

    # Gzip compression
    gzip on;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml;

    # SPA fallback - all routes go to index.html
    location / {
        try_files $uri $uri/ /index.html;
    }

    # Cache static assets
    location /_nuxt/ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    # No cache for index.html
    location = /index.html {
        add_header Cache-Control "no-cache, no-store, must-revalidate";
    }
}
EOF

EXPOSE 80

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["wget", "-q", "--spider", "http://localhost/index.html"] || exit 1

CMD ["nginx", "-g", "daemon off;"]
