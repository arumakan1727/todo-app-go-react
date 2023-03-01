.DEFAULT_GOAL   := help
GO_TEST_FLAG    := -race -shuffle=on
RED      := \033[31m
CYAN     := \033[36m
MAGENTA  := \033[35m
RESET    := \033[0m

.PHONY:	test
test:	## Execute all tests
	go test $(GO_TEST_FLAG) ./...

.PHONY:	test/verbose
test/verbose:	## Execute all tests with verbose output
	go test -v $(GO_TEST_FLAG) ./...

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

.PHONY:	help
help:	## Show Makefile tasks
	@grep -E '^[a-zA-Z_/-]+:.*##' Makefile | \
		awk 'BEGIN {FS = ":.*##\\s*"}; {printf "$(CYAN)%-24s$(RESET) %s\n", $$1, $$2}'
