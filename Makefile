GO_TEST_FLAG    := -race -shuffle=on
RED       := \033[31m
CYAN      := \033[36m
MAGENTA   := \033[35m
RESET     := \033[0m
DB_SERVICE      := postgres
REDIS_SERVICE   := redis


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
	docker compose down --rmi all --volumes --remove-orphans

.PHONY:	db/client
db/client:	## Launch db client
	docker compose exec $(DB_SERVICE) psql -U todoapp todoapp

.PHONY:	db/bash
db/bash:	 ## Launch bash in db service
	docker compose exec $(DB_SERVICE) bash

.PHONY:	help
help:	## Show tasks
	@grep -E '^[a-zA-Z_/-]+:' Makefile | \
		awk 'BEGIN {FS = ":(.*##\\s*)?"}; {printf "$(CYAN)%-24s$(RESET) %s\n", $$1, $$2}'

.PHONY:	help/docker-compose-vars
help/docker-compose-vars:	## Show variables in docker-compose.yml
	./scripts/list-env-params.rb docker-compose.yml
