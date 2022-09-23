.ONESHELL:
SHELL = /bin/bash
MIGRATIONS_FOLDER=$(PWD)/migrations
DB_FILE=$(PWD)/db.sqlite
MAKEFLAGS += --no-print-directory
UI_DIR = "manga-reader-ui"

init:
	@rm -rf vendor
	go mod tidy
	go mod vendor -v
	cd $(UI_DIR)
	npm i

flyway-run:
	@docker run \
		-v $(MIGRATIONS_FOLDER):/flyway/sql \
		-v $(DB_FILE):/flyway/db \
		--network host flyway/flyway:latest-alpine \
		-url=jdbc:sqlite:/flyway/db $(FLYWAY_OPTS) $(FLYWAY_CMD)

db-migrate: FLYWAY_CMD=migrate
db-migrate: flyway-run

ui-build:
	@cd $(UI_DIR)
	npm run build

run:
	@DB_FILE_PATH=db.sqlite \
	WEB_PORT=5678 \
	WEB_CORS=true \
	WEB_PROXY=true \
	DB_DEBUG=false \
	go run main.go

.SILENT: test-start test-stop test-db-init
test-start:
	@docker-compose --profile test up --build --detach

test-stop:
	@docker-compose --profile test down
	@docker-compose --profile test rm -f

test-db-init: DB_FILE := "$(PWD)/tests/helpers/repo/test.db"
test-db-init:
	make flyway-run DB_FILE=$(DB_FILE) FLYWAY_OPTS=-cleanDisabled="false" FLYWAY_CMD=clean >/dev/null 2>&1
	make flyway-run DB_FILE=$(DB_FILE) FLYWAY_CMD=migrate >/dev/null 2>&1

test: test-stop test-db-init test-start
	go test -v ./tests...

stop: test-stop