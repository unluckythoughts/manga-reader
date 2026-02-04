# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add git gcc musl-dev sqlite-dev

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies with cache mount
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

# Copy source code
COPY . ./

# Build the application with cache mount
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o book-reader .

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add ca-certificates sqlite-libs

# Set working directory
WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/book-reader .

# Copy static files if they exist
# COPY --from=builder /app/public ./public

# Create directory for SQLite database
RUN mkdir -p /root/data

# Expose port
EXPOSE 8080

# Run the application
CMD ["./book-reader"]
