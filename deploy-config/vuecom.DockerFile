# Dockerfile for Vue.js application
FROM node:24-alpine AS vue-builder
# FROM oven/bun:debian AS vue-builder

RUN corepack enable
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
WORKDIR /app/frontend
COPY ./frontend/pnpm-lock.yaml ./
COPY ./frontend/package*.json ./
RUN --mount=type=cache,target=/pnpm/store \
    pnpm fetch
COPY ./frontend/package.json ./
RUN --mount=type=cache,target=/pnpm/store \
    pnpm install --frozen-lockfile --prod --offline

# RUN npm install    
COPY ./frontend/ ./
# RUN ls && npm run build
# RUN npm run build
RUN pnpm build


FROM golang:1.25-alpine AS go-builder

WORKDIR /app/backend

COPY ./backend ./
# RUN cd services/gateway_service && go mod download

# COPY ./backend/services/gateway_service/ ./
RUN cd services/gateway_service && ls && go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./bin/server ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=go-builder /app/backend/services/gateway_service/bin/server ./

COPY --from=vue-builder /app/frontend/dist ./dist

# EXPOSE 2500

CMD ["./server"]

