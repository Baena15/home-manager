# ─── Production Dockerfile ──────────────────────────────────────────
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install git and ca-certificates for fetching dependencies
RUN apk add --no-cache git ca-certificates

# Copy dependency files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bin/api ./cmd/api

# ─── Runtime stage ──────────────────────────────────────────────────
FROM alpine:3.19

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk add --no-cache ca-certificates

# Copy binary and static assets
COPY --from=builder /app/bin/api /app/api
COPY --from=builder /app/web /app/web
COPY --from=builder /app/internal/store/migrations /app/migrations

# Railway sets PORT automatically
EXPOSE 8080

CMD ["/app/api"]
