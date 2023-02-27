.DEFAULT_GOAL   := help
GO_TEST_FLAG    := -race -shuffle=on
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

.PHONY:	db/client
db/client:	## Launch db client
	docker compose exec $(DB_SERVICE) psql -U todoapp todoapp

.PHONY:	db/bash
db/bash:	 ## Launch bash in db service
	docker compose exec $(DB_SERVICE) bash

.PHONY:	db/migration/status
db/migration/status:	## Show migration status
	dbmate status

.PHONY:	db/migration/gen
db/migration/gen:	## Generate migration file from 'db/schema.sql'
	@./scripts/make-migration.sh

.PHONY:	db/migration/up
db/migration/up:	## Apply migration
	dbmate up

.PHONY:	db/migration/dump
db/migration/dump:	## Dump table schema
	dbmate dump

.PHONY:	lint/spectral
lint/spectral:	## Lint swagger.yaml
	spectral lint --ruleset .spectral.yaml -f pretty api-spec/swagger.yaml

.PHONY:	help
help:	## Show tasks
	@grep -E '^[a-zA-Z_/-]+:' Makefile | \
		awk 'BEGIN {FS = ":(.*##\\s*)?"}; {printf "$(CYAN)%-24s$(RESET) %s\n", $$1, $$2}'

.PHONY:	help/docker-compose-vars
help/docker-compose-vars:	## Show variables in docker-compose.yml
	./scripts/list-env-params.rb docker-compose.yml
