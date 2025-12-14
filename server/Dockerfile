# ------------------------------------------------------
# 1. Build Stage Server
# ------------------------------------------------------
FROM golang:1.25-alpine AS go_builder

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source files
COPY . .

# Build with CGO enabled for SQLite
RUN GOOS=linux go build -o /app/server_bin

# ------------------------------------------------------
# 2. Runtime Stage
# ------------------------------------------------------
FROM alpine:3.20

WORKDIR /app

# Copy binary and required files
COPY --from=go_builder /app/server_bin ./server
COPY --from=go_builder /app/config ./config

EXPOSE 8080

CMD ["./server"]
