## help: prints help for targets with comments
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'
confirm:
	@echo -n 'Are you sure? [y/N]: ' && read ans && [ $${ans:-N} == y ]

## api/run: run the cmd/api application
api/run:
	go run ./cmd/api -db-dsn=postgres://greenlight:1234@localhost/greenlight?sslmode=disable

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
