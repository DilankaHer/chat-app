# ------------------------------------------------------
# 1. Build Stage Server
# ------------------------------------------------------
FROM golang:1.25-alpine AS go_builder

WORKDIR /app/server

# Cache dependencies
COPY server/go.mod server/go.sum ./
RUN go mod download

# Copy source files
COPY server .

# Build with CGO enabled for SQLite
RUN GOOS=linux go build -o /app/server_bin

# ------------------------------------------------------
# 2. Build Stage Client
# ------------------------------------------------------
FROM node:24-alpine as vite_builder

WORKDIR /app/client

# Cache dependencies
COPY client/package.json client/package-lock.json ./
RUN npm install

# Copy source files
COPY client .

# Build with CGO enabled for SQLite
RUN npm run build

# ------------------------------------------------------
# 2. Runtime Stage
# ------------------------------------------------------
FROM alpine:3.20

WORKDIR /app

# Copy binary and required files
COPY --from=go_builder /app/server_bin ./server_bin
COPY --from=go_builder /app/server/config ./config
COPY --from=vite_builder /app/client/dist ./static

EXPOSE 8080

CMD ["./server_bin"]
