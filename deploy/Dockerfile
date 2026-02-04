
# Migrations stage
FROM flyway/flyway:latest-alpine AS migrations

COPY migrations /flyway/sql

RUN apk add sqlite

# Create directory for SQLite database
RUN mkdir -p /flyway/data

# Run migrations on SQLite database
RUN flyway \
    -url=jdbc:sqlite:/flyway/data/book-reader.db \
    -user=dummy \
    -password=dummy \
    -locations=filesystem:/flyway/sql \
    migrate 2>&1 | tee ./migration.log

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

# Copy the migrated SQLite database from migrations stage
COPY --from=migrations /flyway/data/book-reader.db /root/data/book_reader.db
COPY --from=migrations /flyway/migration.log /root/migration.log

# Expose port
EXPOSE 8080

# Run the application
CMD ["./book-reader"]
