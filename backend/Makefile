.DEFAULT_GOAL   := help
DB_SERVICE      := postgres
REDIS_SERVICE   := redis
GO_TEST_FLAG    := -race -shuffle=on
RED      := \033[31m
CYAN     := \033[36m
MAGENTA  := \033[35m
RESET    := \033[0m

.PHONY:	test/all
test/all:	## Execute all tests
	go test $(GO_TEST_FLAG) ./...

.PHONY:	test/repository
test/repository:	## Execute tests for ./repository/...
	go test $(GO_TEST_FLAG) ./repository/...

.PHONY:	test/domain
test/domain:	## Execute tests for ./domain/...
	go test $(GO_TEST_FLAG) ./domain/...

.PHONY:	test/presenter
test/presenter:	## Execute tests for ./presenter/...
	go test $(GO_TEST_FLAG) ./presenter/...

.PHONY:	fmt
fmt:	## Format and overwrite files
	go fmt ./...

.PHONY:	gen/from-swagger
gen/from-swagger:	## Generate codes from swagger.yaml
	go generate ./presenter/restapi/

.PHONY:	gen/from-pgsql-schema
gen/from-pgsql-schema:	## Generate codes from swagger.yaml
	go generate ./repository/pgsql/sqlcgen/
	./_scripts/after-sqlc-generate.sh

.PHONY:	lint/go
lint/go:	## Lint golang source codes
	golangci-lint run ./...


DB_USER ?= todouser
DB_HOST ?= 127.0.0.1
DB_PORT ?= 25432
MAIN_DB_NAME ?= tododb
TEST_DB_NAME ?= $(MAIN_DB_NAME)__test

# dbmate コマンドは環境変数からデータベースURLを指定できるのでそれを使う
export MAIN_DATABASE_URL := postgres://$(DB_USER)@$(DB_HOST):$(DB_PORT)/$(MAIN_DB_NAME)?sslmode=disable
export TEST_DATABASE_URL := postgres://$(DB_USER)@$(DB_HOST):$(DB_PORT)/$(TEST_DB_NAME)?sslmode=disable
export DBMATE_MIGRATIONS_DIR := _postgres/migrations

.PHONY:	db/client/maindb
db/client/maindb:	## Launch postgres client for main database
	docker compose exec $(DB_SERVICE) psql -U $(DB_USER) $(MAIN_DB_NAME)

.PHONY:	db/client/testdb
db/client/testdb:	## Launch postgres client for test database
	docker compose exec $(DB_SERVICE) psql -U $(DB_USER) $(TEST_DB_NAME)

.PHONY:	db/bash
db/bash:	 ## Launch bash in db service
	docker compose exec $(DB_SERVICE) bash

.PHONY:	db/urls
db/urls:	## Show database URLs
	@echo "MAIN_DATABASE_URL='$(MAIN_DATABASE_URL)'"
	@echo "TEST_DATABASE_URL='$(TEST_DATABASE_URL)'"

.PHONY:	db/migration/status
db/migration/status:	## Show migration status
	@echo "$(CYAN)### Status of '$(MAIN_DB_NAME)':$(RESET)"
	dbmate --env=MAIN_DATABASE_URL status
	@echo "\n$(CYAN)### Status of '$(TEST_DB_NAME)':$(RESET)"
	dbmate --env=TEST_DATABASE_URL status

.PHONY:	db/migration/new
db/migration/new:	## Create migration file
	@./scripts/new-migration.sh

.PHONY:	db/migration/up
db/migration/up:	## Apply migration
	@echo "$(CYAN)### Applying to '$(TEST_DB_NAME)':$(RESET)"
	dbmate --env=TEST_DATABASE_URL up
	@echo "\n$(CYAN)### Applying to '$(MAIN_DB_NAME)':$(RESET)"
	dbmate --env=MAIN_DATABASE_URL up

.PHONY:	__db/migration/down
__db/migration/down:	## Revert migration
	@echo "$(CYAN)### Reverting '$(TEST_DB_NAME)':$(RESET)"
	dbmate --env=TEST_DATABASE_URL down
	@echo "\n$(CYAN)### Reverting '$(MAIN_DB_NAME)':$(RESET)"
	dbmate --env=MAIN_DATABASE_URL down

.PHONY:	db/migration/dump
db/migration/dump:	## Dump table schema
	dbmate --env=MAIN_DATABASE_URL dump

.PHONY:	help
help:	## Show Makefile tasks
	@grep -E '^[a-zA-Z_/-]+:.*##' Makefile | \
		awk 'BEGIN {FS = ":.*##\\s*"}; {printf "$(CYAN)%-24s$(RESET) %s\n", $$1, $$2}'
