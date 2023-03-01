.DEFAULT_GOAL   := help
DB_SERVICE      := postgres
REDIS_SERVICE   := redis
RED      := \033[31m
CYAN     := \033[36m
MAGENTA  := \033[35m
RESET    := \033[0m


.PHONY:	docker/build/dev
docker/build/dev:	## Build docker image to local development
	docker compose build

.PHONY:	docker/up
docker/up:	## Do docker compose up
	docker compose up -d

.PHONY:	docker/down
docker/down:	## Do docker compose down
	docker compose down

.PHONY:	docker/restart
docker/restart:	## Do docker compose restart
	docker compose restart

.PHONY:	docker/logs
docker/logs:	## Tail docker composee logs
	docker compose logs -f

.PHONY:	docker/ps
docker/ps:	## Check container status
	docker compose ps

.PHONY:	__docker/down-remove
__docker/down-remove:	## Remove containers and volumes
	docker compose down --rmi local --volumes --remove-orphans

DB_USER ?= todouser
DB_HOST ?= 127.0.0.1
DB_PORT ?= 25432
DB_NAME ?= tododb
# dbmate コマンドは環境変数からデータベースURLを指定できるのでそれを使う
export MAIN_DATABASE_URL := postgres://$(DB_USER)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
export TEST_DATABASE_URL := postgres://$(DB_USER)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)__test?sslmode=disable

.PHONY:	db/client
db/client:	## Launch db client
	docker compose exec $(DB_SERVICE) psql -U $(DB_USER) $(DB_NAME)

.PHONY:	db/bash
db/bash:	 ## Launch bash in db service
	docker compose exec $(DB_SERVICE) bash

.PHONY:	db/urls
db/urls:	## Show database URLs
	@echo "MAIN_DATABASE_URL='$(MAIN_DATABASE_URL)'"
	@echo "TEST_DATABASE_URL='$(TEST_DATABASE_URL)'"

.PHONY:	db/migration/status
db/migration/status:	## Show migration status
	@echo "$(CYAN)### Status of '$(DB_NAME)':$(RESET)"
	dbmate --env=MAIN_DATABASE_URL status
	@echo "\n$(CYAN)### Status of '$(DB_NAME)__test':$(RESET)"
	dbmate --env=TEST_DATABASE_URL status

.PHONY:	db/migration/new
db/migration/new:	## Generate migration file from 'db/schema.sql'
	@./scripts/new-migration.sh

.PHONY:	db/migration/up
db/migration/up:	## Apply migration
	@echo "$(CYAN)### Applying to '$(DB_NAME)__test':$(RESET)"
	dbmate --env=TEST_DATABASE_URL up
	@echo "\n$(CYAN)### Applying to '$(DB_NAME)':$(RESET)"
	dbmate --env=MAIN_DATABASE_URL up

.PHONY:	__db/migration/down
__db/migration/down:	## Revert migration
	@echo "$(CYAN)### Reverting '$(DB_NAME)__test':$(RESET)"
	dbmate --env=TEST_DATABASE_URL down
	@echo "\n$(CYAN)### Reverting '$(DB_NAME)':$(RESET)"
	dbmate --env=MAIN_DATABASE_URL down

.PHONY:	db/migration/dump
db/migration/dump:	## Dump table schema
	dbmate --env=MAIN_DATABASE_URL dump

.PHONY:	lint/spectral
lint/spectral:	## Lint swagger.yaml
	spectral lint --ruleset .spectral.yaml -f pretty api-spec/swagger.yaml

.PHONY:	help
help:	## Show Makefile tasks
	@grep -E '^[a-zA-Z_/-]+:' Makefile | \
		awk 'BEGIN {FS = ":(.*##\\s*)?"}; {printf "$(CYAN)%-24s$(RESET) %s\n", $$1, $$2}'

.PHONY:	help/docker-compose-vars
help/docker-compose-vars:	## Show variables in docker-compose.yml
	./scripts/list-env-params.rb docker-compose.yml
