# Dockerfile for Vue.js application
FROM node:24-alpine AS vue-builder
# FROM oven/bun:debian AS vue-builder

RUN corepack enable
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"

WORKDIR /app/frontend

COPY ./frontend/pnpm-lock.yaml ./
COPY ./frontend/package*.json ./

# RUN --mount=type=cache,target=/pnpm/store \
#     pnpm fetch
COPY ./frontend/package.json ./

RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile

# RUN npm install    
COPY ./frontend/ ./
# RUN ls && npm run build
# RUN npm run build
RUN pnpm build


FROM golang:1.25-alpine AS go-builder

WORKDIR /app/backend

COPY ./backend/services/analytics_service/go* ./services/analytics_service/
COPY ./backend/services/inventory_service/go* ./services/inventory_service/
COPY ./backend/services/payment_service/go* ./services/payment_service/
COPY ./backend/services/catalog_service/go* ./services/catalog_service/
COPY ./backend/services/gateway_service/go* ./services/gateway_service/

COPY ./backend/shared ./shared


RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    cd services/gateway_service && go mod download
COPY ./backend ./


RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \ 
    cd services/gateway_service && CGO_ENABLED=0 GOOS=linux go build -o ./bin/server ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=go-builder /app/backend/services/gateway_service/bin/server ./

COPY --from=vue-builder /app/frontend/dist ./dist

USER nonroot:nonroot

# EXPOSE 2500

CMD ["./server"]

