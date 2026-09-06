# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Compile optimized static binary without debug symbols
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/server ./cmd/api

# Run stage
FROM alpine:3.20

WORKDIR /app

# Install certificates and timezone data for HTTPS and timestamps
RUN apk --no-cache add ca-certificates tzdata

# Create dedicated non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Create and set permissions for uploads directory
RUN mkdir -p /app/uploads && chown -R appuser:appgroup /app

# Copy binary from builder
COPY --from=builder --chown=appuser:appgroup /app/server .

USER appuser

EXPOSE 8090

# Native healthcheck using wget
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8090/health || exit 1

CMD ["./server"]
