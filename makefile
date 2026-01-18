.ONESHELL:
SHELL = powershell.exe
.SHELLFLAGS = -NoProfile -Command
MIGRATIONS_FOLDER=$(PWD)/migrations
DB_FILE=$(PWD)/db.sqlite
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

# Build the application locally
build:
	go build -o book-reader.exe .

# Docker build
docker-build:
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# Docker build and run
docker-up: docker-build
	docker-compose -f deploy/docker-compose.yml up -d

# Docker stop and remove
docker-down:
	docker-compose -f deploy/docker-compose.yml down -v

# Run tests
test:
	go test ./tests/... -v

# Run API tests
test-api:
	go test ./tests/api/... -v

# Run tests with coverage
test-coverage:
	go test ./tests/api/... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

flyway-run:
	@docker run \
		-v $(MIGRATIONS_FOLDER):/flyway/sql \
		-v $(DB_FILE):/flyway/db \
		--network host flyway/flyway:latest-alpine \
		-url=jdbc:sqlite:/flyway/db $(FLYWAY_OPTS) $(FLYWAY_CMD)

db-migrate: FLYWAY_CMD=migrate
db-migrate: flyway-run

# Clean build artifacts
clean:
	Remove-Item -ErrorAction SilentlyContinue book-reader.exe
	Remove-Item -ErrorAction SilentlyContinue coverage.out
	Remove-Item -ErrorAction SilentlyContinue coverage.html
