# =============================================================================
# CortexOps Frontend — Multi-Stage Docker Build
# =============================================================================
# References:
#   https://nextjs.org/docs/app/building-your-application/deploying#docker-image
#   Audit finding: TD-002
#
# Build stages:
#   1. deps    — install production + dev dependencies (layer-cached)
#   2. builder — compile Next.js application (standalone output mode)
#   3. runner  — minimal production image (~80% smaller than builder)
#
# Required next.config.ts setting: output: 'standalone'
# =============================================================================

# =============================================================================
# Stage 1: deps
# Install dependencies only. This layer is cached independently from the build
# so that npm install does not re-run on every source file change.
# =============================================================================
FROM node:20-alpine AS deps

# Add libc6-compat for native module compatibility on Alpine.
RUN apk add --no-cache libc6-compat

WORKDIR /app

# Copy only the package manifests first to maximise cache hit rate.
# node_modules is rebuilt only when package.json or package-lock.json changes.
COPY package.json package-lock.json ./

# Use `npm ci` (clean install) for reproducible, lockfile-exact installs.
# This is the correct production-grade install command (not `npm install`).
RUN npm ci

# =============================================================================
# Stage 2: builder
# Build the Next.js application. Depends on the deps stage so that node_modules
# is available during compilation.
# =============================================================================
FROM node:20-alpine AS builder

WORKDIR /app

# Copy installed node_modules from the deps stage.
COPY --from=deps /app/node_modules ./node_modules

# Copy the rest of the application source.
COPY . .

# Copy the playbooks JSON that the build pipeline may need.
# The copy-playbooks.mjs script runs before the build to inline playbook content.
# If the script fails (e.g., docs/playbooks does not exist), the build continues.
RUN node copy-playbooks.mjs || echo "No playbooks found, continuing build."

# Build environment variables.
# NEXT_TELEMETRY_DISABLED=1: Opt out of Next.js anonymous telemetry in CI/CD.
ENV NEXT_TELEMETRY_DISABLED=1 \
    NODE_ENV=production

# Build the Next.js application.
# The 'standalone' output mode (set in next.config.ts) produces a minimal
# server at .next/standalone/ that does not require node_modules at runtime.
RUN npm run build

# =============================================================================
# Stage 3: runner
# Minimal production image. Only the compiled output and static assets are
# copied from the builder stage. node_modules is NOT included.
# =============================================================================
FROM node:20-alpine AS runner

WORKDIR /app

# Security: Run the application as a non-root user.
# The 'nextjs' user and group are created here with explicit IDs to ensure
# consistency across environments and to satisfy Kubernetes runAsNonRoot policy.
RUN addgroup --system --gid 1001 nodejs \
 && adduser  --system --uid  1001 nextjs

ENV NODE_ENV=production \
    NEXT_TELEMETRY_DISABLED=1 \
    # The Next.js standalone server listens on this port.
    PORT=3000 \
    # Bind to all interfaces so the container is reachable from the Docker network.
    HOSTNAME=0.0.0.0

# Copy the pre-built standalone server from the builder stage.
# The standalone directory contains a minimal Node.js server (server.js) and
# only the node_modules required at runtime — nothing more.
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./

# Copy the static asset directory (CSS, JS chunks, images).
# These are served directly by the Next.js standalone server.
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

# Copy the public directory (favicon, OG images, manifests, etc.).
COPY --from=builder --chown=nextjs:nodejs /app/public ./public

# Switch to the non-root user.
USER nextjs

EXPOSE 3000

# Health check: verify the server responds before marking the container healthy.
# Used by Docker Compose 'condition: service_healthy' and Kubernetes readiness probes.
HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:3000/ || exit 1

# Start the standalone Next.js server.
CMD ["node", "server.js"]
