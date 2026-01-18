
# Migrations stage
FROM redgate/flyway:latest-alpine AS migrations

COPY ../migrations /flyway/sql

RUN apk --no-cache add sqlite

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
RUN apk add --no-cache git gcc musl-dev sqlite-dev

# Set working directory
WORKDIR /app

# Copy go mod files
COPY ../go.mod ../go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY ../ ./

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o book-reader .

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates sqlite-libs

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
