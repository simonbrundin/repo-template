# =============================================================================
# Nuxt + Bun Dockerfile
# En enda Dockerfile med named stages för dev och production
#
# Usage:
#   Development:  docker build --target dev -t nuxt-dev . && docker run -it nuxt-dev
#   Production:  docker build -t nuxt-prod . && docker run -p 3000:3000 nuxt-prod
# =============================================================================

# ---- Base: gemensamma dependencies ----
FROM oven/bun:1-alpine AS base

WORKDIR /app

COPY package.json bun.lock ./
RUN bun install --frozen-lockfile

# ---- Dev: fullständig utvecklingsmiljö ----
FROM base AS dev

# Kopiera hela source
COPY . .

# Standard-kommando för dev
CMD ["bun", "run", "dev", "--", "--host", "0.0.0.0"]

# ---- Builder: bygger production-app ----
FROM base AS builder

# Kopiera source och bygg
COPY . .
RUN bun run build

# ---- Production: minimal runtime ----
FROM oven/bun:1-alpine AS production

WORKDIR /app

ENV NODE_ENV=production

# Kopiera byggda filer från builder
COPY --from=builder /app/.output ./.output
COPY --from=builder /app/.nuxt ./.nuxt
COPY --from=builder /app/node_modules ./node_modules

COPY nuxt.config.ts ./

EXPOSE 3000
USER node

CMD ["bun", "run", ".output/server/index.mjs"]
