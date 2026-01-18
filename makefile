.ONESHELL:
SHELL = powershell.exe
.SHELLFLAGS = -NoProfile -Command
MIGRATIONS_FOLDER=$(PWD)/migrations
DB_FILE=$(PWD)/db.sqlite
MAKEFLAGS += --no-print-directory
UI_DIR = manga-reader-ui

# 	if (Test-Path "$(UI_DIR)") { Set-Location "$(UI_DIR)"; npm i }
init:
	go clean -modcache
	go mod tidy
	go mod download

flyway-run:
	@docker run \
		-v $(MIGRATIONS_FOLDER):/flyway/sql \
		-v $(DB_FILE):/flyway/db \
		--network host flyway/flyway:latest-alpine \
		-url=jdbc:sqlite:/flyway/db $(FLYWAY_OPTS) $(FLYWAY_CMD)

db-migrate: FLYWAY_CMD=migrate
db-migrate: flyway-run
