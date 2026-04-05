# Frontend stage
FROM node:18-alpine AS frontend

# Set working directory
WORKDIR /app

# Copy package.json and package-lock.json
COPY app/package.json app/package-lock.json ./

# Install dependencies with cache mount
RUN --mount=type=cache,target=/root/.npm \
		npm install

# Copy frontend source code
COPY app/ ./

# Build the frontend
RUN npm run build

# Backend stage
FROM golang:1.26-alpine AS backend

# Install build dependencies
RUN apk add --no-cache \
    git \
    gcc \
    musl-dev \
    sqlite \
    sqlite-dev

# Set working directory
WORKDIR /app

# Copy book-reader go mod files
COPY go.mod go.sum ./

# Download dependencies with cache mount
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

# Copy book-reader source code
COPY . .

# COPY static files from frontend stage
COPY --from=frontend /app/dist ./public

# Build the application with cache mount
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=1 GOOS=linux go build -tags musl -o book-reader .

# Runtime stage
FROM alpine:latest AS runtime

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    wget \
    sqlite-libs

# Set working directory
WORKDIR /app

# Copy the binary from backend
COPY --from=backend /app/book-reader .

# Copy static files if they exist
# COPY --from=frontend /app/dist ./public

# Create directory for SQLite database
RUN mkdir -p /app/data

# Expose port
EXPOSE 8080

# Run the application
CMD ["./book-reader"]
