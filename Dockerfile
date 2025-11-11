# Multi-stage build for Go Agent Template
# Stage 1: Build
FROM golang:1.21-alpine AS builder

LABEL maintainer="Mindra Platform"
LABEL description="Go Agent Template for Mindra Platform"

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod ./

# Download dependencies (if any)
RUN go mod download || true

# Copy source code
COPY . .

# Build the binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o agent \
    main.go

# Stage 2: Production runtime
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk --no-cache add ca-certificates curl

# Create non-root user
RUN addgroup -g 1001 -S gouser && \
    adduser -S gouser -u 1001 -G gouser && \
    chown -R gouser:gouser /app

# Copy binary from builder
COPY --from=builder /app/agent .

USER gouser

# Expose port
EXPOSE 8003

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8003/health || exit 1

# Run agent
CMD ["./agent"]
