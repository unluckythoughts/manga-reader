.ONESHELL:
SHELL = powershell.exe
.SHELLFLAGS = -NoProfile -Command
MIGRATIONS_FOLDER=$(PWD)/migrations
DB_FILE=$(PWD)/data/book_reader.db
MAKEFLAGS += --no-print-directory
UI_DIR = book-reader-ui

# Docker variables
DOCKER_IMAGE=book-reader
DOCKER_TAG=latest
ENV_FILE=deploy/.env

# 	if (Test-Path "$(UI_DIR)") { Set-Location "$(UI_DIR)"; npm i }
init:
	go clean -modcache
	go mod tidy
	go mod download

build:
	docker build -f deploy/Dockerfile -t book-reader:latest ..

start: db-migrate
	docker-compose -f deploy/docker-compose.yml up -d

stop:
	docker-compose -f deploy/docker-compose.yml down -v

run: start

# Docker stop and remove
docker-down:
	docker-compose -f deploy/docker-compose.yml down -v

# Run tests
test:
	go test ./tests/... -v

# Run API tests
test-api:
	go test ./tests/api/... -v

# Run connector integration tests
test-connector:
	go test ./tests/connector/... -v -timeout 5m

# Run connector tests in short mode (skip network requests)
test-connector-short:
	go test ./tests/connector/... -v -short

# Run connector benchmarks
test-connector-bench:
	go test ./tests/connector/... -bench=. -benchmem

# Run tests with coverage
test-coverage:
	go test ./tests/api/... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

flyway-run:
	docker run --rm -v "$(CURDIR)/migrations:/flyway/sql:ro" -v "$(CURDIR)/data:/flyway/data" -w /flyway flyway/flyway:latest-alpine -url=jdbc:sqlite:/flyway/data/book_reader.db $(FLYWAY_OPTS) $(FLYWAY_CMD)
	docker run --rm -v "$(CURDIR)/migrations:/flyway/sql:ro" -v "$(CURDIR)/data:/flyway/data" -w /flyway flyway/flyway:latest-alpine -url=jdbc:sqlite:/flyway/data/book_reader_auth.db $(FLYWAY_OPTS) $(FLYWAY_CMD)

db-migrate: FLYWAY_CMD=migrate
db-migrate: flyway-run

# Clean build artifacts
clean:
	Remove-Item -ErrorAction SilentlyContinue book-reader.exe
	Remove-Item -ErrorAction SilentlyContinue coverage.out
	Remove-Item -ErrorAction SilentlyContinue coverage.html
