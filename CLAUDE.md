# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What This Is

VueCom is an open-source, self-hostable e-commerce platform (Shopify/Magento alternative) built on MACH principles. It is in early development. The frontend is a Vue 3 SPA and the backend is a Go modular monolith exposed through a single Fiber HTTP gateway.

## Commands

### Root (pnpm, run from repo root)

```bash
pnpm dev              # Start frontend (Vite HMR) + backend (air hot-reload) concurrently
pnpm build            # Build Go binary + Vue app → .output/gateway/
pnpm start            # Run production binary from .output/
pnpm install:server   # go mod download for all backend services
```

### Frontend (run from `frontend/`)

```bash
pnpm lint             # Run oxlint + eslint with auto-fix
pnpm lint:oxlint      # oxlint only
pnpm lint:eslint      # eslint only
pnpm format           # Prettier
pnpm type-check       # vue-tsc --build
pnpm test:unit        # Vitest (jsdom)
pnpm test:e2e         # Playwright
```

### Backend (run from `backend/services/gateway/`)

```bash
make build    # go build -o bin/server ./cmd/server
make prod     # Stripped production build
go test ./... # Run tests (no Makefile target yet)
```

### Infrastructure

```bash
docker-compose up vuecom-db redis rabbitmq   # Dev infra only
docker-compose up                            # Full stack (app on host :3500)
```

## Environment Setup

Copy and fill `backend/services/gateway/.env` (see `.env.example`). Required vars: `APP_PG_HOST`, `APP_PG_USER`, `APP_PG_PASSWORD`, `GATE_PG_DBNAME`, `APP_PG_PORT`, `APP_REDIS_URL`, `CLOUDINARY_KEY`, `CLOUDINARY_SECRET`, `CLOUDINARY_CLOUD_NAME`, `SECRET_KEY`, `API_ENC_KEY`. Also create `deploy-config/redis/secrets.conf` from its example.

Dev tooling prerequisites: Node.js ≥ 20.10, pnpm, Go ≥ 1.22.1, Docker, and `air` (`go install github.com/air-verse/air@latest`).

## Architecture

### Single Binary, Internal gRPC

All domain services (catalog, orders, inventory, payment, notification, analytics) compile into one Go binary alongside the `gateway` HTTP server. They communicate via **in-process gRPC over `bufconn`** (an in-memory listener), not real network sockets. This gives typed proto contracts and service isolation without network overhead. Setup lives in `backend/services/gateway/internal/grpc/`.

### Proto Contracts

All `.proto` definitions and their generated Go code live in `backend/shared/proto/go/{catalog,orders,inventory}/`. When you add or modify a service, regenerate from there.

### Go Workspace

`backend/go.work` links all service modules and `shared/` together. Each service (`catalog`, `orders`, etc.) is its own Go module that imports `shared` via the workspace replace directive — no need to publish or version it.

### Frontend ↔ Backend

In dev, Vite proxies `/api/*` → `localhost:2500`. In production, the Go binary serves the Vue `dist/` as a SPA fallback and exposes the API at `/api/v1/`.

### Event Bus

Async cross-service messaging uses **RabbitMQ** (amqp091-go). Queue names and event type constants are defined centrally in `backend/shared/events/events.go`. Catalog has the reference pub/sub implementation in `backend/services/catalog/internal/pubsub/`.

### API Surface

- `GET /api/health`
- `/api/v1/backend/...` — Admin routes: JWT cookie (`backend_access_token`) + RBAC required, Redis rate-limited at 1000 req/min
- `/api/v1/customer/...` — Storefront routes: rate-limited at 100 req/min per customer
- `GET *` — SPA fallback serving `dist/index.html`

### Auth Model

JWT access token (15 min, HTTP-only cookie) + refresh token (7 days, hashed in DB) with device-ID cookie for session binding. Admin registration is invite-only: an admin creates a signup token → user receives a URL → completes registration via HOTP code validation.

### RBAC

Custom `rbac.PermissionSet` stored in an LRU cache keyed by user ID (`backend/services/gateway/internal/global/`). The `HasPermission(...)` middleware in `api/v1/middlewares/rbac.go` guards backend routes.

### Repository Pattern

Each domain area has a dedicated GORM repository (e.g. `backendUserRepository`, `productRepository`) accessed through the `GormPGDatabase` interface defined in `backend/services/gateway/internal/types/`. Dependency injection flows through the `Deps` struct.

## Key Tech

| Layer | Choice |
|---|---|
| Frontend | Vue 3 + Pinia + Vue Router 4 + TailwindCSS v4 + reka-ui |
| Build | rolldown-vite (Rolldown-based Vite fork) |
| Backend | Go + Fiber v2 (fasthttp) |
| ORM | GORM v2 + pgx v5 (PostgreSQL 18) |
| Cache | Redis (go-redis v9) |
| Queue | RabbitMQ (amqp091-go) |
| gRPC | google.golang.org/grpc v1.80 via bufconn |
| Logging | Zap (uber-go/zap) |
| Media | Cloudinary |
