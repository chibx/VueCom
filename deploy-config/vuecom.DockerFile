# Dockerfile for Vue.js application
# FROM node:24-alpine AS vue-builder
FROM oven/bun:debian AS vue-builder

WORKDIR /app/frontend

COPY ./frontend/package*.json ./
RUN bun install

COPY ./frontend/ ./

RUN bun run build

FROM golang:1.25-alpine AS go-builder

WORKDIR /app/backend

COPY ./backend/services/gateway_service/go.mod ./backend/services/gateway_service/go.sum ./
RUN go mod download

COPY ./backend/services/gateway_service/ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o bin/server .

FROM gcr.io/distroless/static-debian12:nonroot

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=go-builder /app/backend/bin/server ./

COPY --from=vue-builder /app/frontend/dist ./dist

# EXPOSE 2500

CMD ["./server"]

