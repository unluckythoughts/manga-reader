.ONESHELL:
MAKEFLAGS += --no-print-directory

ifeq ($(OS),Windows_NT)
	SHELL = powershell.exe
	.SHELLFLAGS = -NoProfile -Command
	BINARY = book-reader.exe
	RM = Remove-Item -ErrorAction SilentlyContinue
	RUN_ARGS = $$env:LOG_LEVEL="debug"; $$env:AUTH_JWT_KEY="debug-secret-key-change-in-production"; 
else
	SHELL = /bin/bash
	.SHELLFLAGS = -e -c
	BINARY = book-reader
	RM = rm -f
	RUN_ARGS = LOG_LEVEL=debug AUTH_JWT_KEY=debug-secret-key-change-in-production 
endif

DOCKER_COMPOSE_FILE = deploy/docker-compose.yml
DOCKER_COMPOSE = docker compose
MIGRATIONS_FOLDER=$(PWD)/migrations
DB_FILE=$(PWD)/data/book_reader.db
UI_DIR = book-reader-ui

# Docker variables
DOCKER_IMAGE=book-reader
DOCKER_TAG=latest
ENV_FILE=deploy/.env

init:
	go clean -modcache
	go mod tidy
	go mod download

build:
	docker build -f deploy/Dockerfile -t book-reader:latest ..

# Build locally for debugging
build-local:
	go build -gcflags="all=-N -l" -o $(BINARY) .

# Build frontend into public and compile executable
build-exe:
ifeq ($(OS),Windows_NT)
	if (-not (Test-Path public)) { New-Item -ItemType Directory -Path public | Out-Null }
	pushd app; npm install; npx vite build --outDir ../public; popd
else
	mkdir -p public
	(cd app && npm install && npx vite build --outDir ../public)
endif
	go build -v -o $(BINARY) .

# Run locally with debug output
debug: build-local db-migrate
	$(RUN_ARGS) ./$(BINARY)


start: db-migrate
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) up --build -d

stop:
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) down -v

run: start

# Docker stop and remove
docker-down:
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) down -v

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
	$(RM) $(BINARY), coverage.out, coverage.html
