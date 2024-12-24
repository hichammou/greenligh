# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: prints help for targets with comments
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N]: ' && read ans && [ $${ans:-N} == y ]

# ==================================================================================== 
# DEVELOPMENT
# ==================================================================================== #

## api/run: run the cmd/api application
api/run:
	go run ./cmd/api -db-dsn=postgres://greenlight:1234@localhost:5433/greenlight?sslmode=disable

## db/psql: connect to the databas using psql
db/psql:
	docker exec -it 3b5d psql postgres://greenlight:1234@localhost/greenlight?sslmode=disable

## db/migrations/new name=$1: create a new database migration
db/migrations/new:
	@echo 'Creating migration files for ${name}'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all the up migrations
db/migrations/up: confirm
	@echo 'Running up migrations ...'
	migrate -path ./migrations -database postgres://greenlight:1234@localhost/greenlight?sslmode=disable up

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: tidy dependecies and format, vet and test all code
.PHONY: audit
audit: vendor
	# @echo 'Formatting code...' Disable format for now
	# go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

## vendor: tidy and vendor dependecies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependecies...'
	go mod tidy
	go mod verify
	@echo 'vendoring dependecies...'
	go mod vendor

# ==================================================================================== #
# BUILD
# ==================================================================================== #

current_time = $(shell date --iso-8601=seconds)
linker_flags = '-s -X main.buildTime=${current_time}'

## build/api: build the cmd/api application
.PHONY: api/build
api/build:
	@echo 'Building cmd/api...'
	go build -ldflags=${linker_flags} -o=./bin/api ./cmd/api
	GOOS=linux GOARCH=amd64 go build -ldflags=${linker_flags} -o=./bin/linux_amd64/api ./cmd/api
