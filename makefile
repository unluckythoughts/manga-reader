.ONESHELL:
SHELL = /bin/bash
MIGRATIONS_FOLDER=$(PWD)/migrations
DB_FILE=$(PWD)/db.sqlite


init:
	@rm -rf vendor
	go mod vendor -v

db-migrate:
	@docker run \
		-v $(MIGRATIONS_FOLDER):/flyway/sql \
		-v $(DB_FILE):/flyway/db \
		--network host flyway/flyway:latest-alpine \
		-url=jdbc:sqlite:/flyway/db migrate