# Multi-stage Dockerfile for dago-node-planner

# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
ARG VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.gitCommit=${GIT_COMMIT}" \
    -o /build/node-planner \
    ./cmd/node-planner

# Runtime stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 planner && \
    adduser -D -u 1000 -G planner planner

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/node-planner /app/node-planner

# Copy prompts
COPY --chown=planner:planner prompts /app/prompts

# Note: JSON schemas are embedded in dago-libs binary

# Switch to non-root user
USER planner

# Expose port
EXPOSE 8080

# Set environment variables
ENV PLANNER_SERVER_HOST=0.0.0.0 \
    PLANNER_SERVER_PORT=8080 \
    PLANNER_PROMPT_PATH=/app/prompts \
    PLANNER_LOG_LEVEL=info \
    PLANNER_LOG_FORMAT=json

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the binary
ENTRYPOINT ["/app/node-planner"]
