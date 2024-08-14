name ?= ''
.PHONY: help migrate-up migrate-down fresh migrate-create sqlc
help: ## Print help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
migrate-up: ## run migrate up
	$(info migrate up running....)
#	@migrate -path db/migration -database "postgresql://amer:amer@127.0.0.1:5432/simple_bank?sslmode=disable" force 1
#	@migrate -path db/migration -database "postgresql://amer:amer@127.0.0.1:5432/simple_bank?sslmode=disable" -verbose up
	@goose -dir db/migrations postgres "postgresql://amer:amer@127.0.0.1:5432/simple_bank?sslmode=disable" up
migrate-down: ## run migrate down
	$(info migrate down running....)
#	@migrate -path db/migration -database "postgresql://amer:amer@127.0.0.1:5432/simple_bank?sslmode=disable" force 1
#	@migrate -path db/migration -database "postgresql://amer:amer@127.0.0.1:5432/simple_bank?sslmode=disable" -verbose down
	@goose -dir db/migrations postgres "postgresql://amer:amer@127.0.0.1:5432/simple_bank?sslmode=disable" down

migrate-create: ## Target to create a new migration file using goose
	@if [ -z "${name}" ]; then \
		echo "Error: NAME is not set."; \
		echo "Usage: make migrate-create NAME=<your_migration_name>"; \
		exit 1; \
	fi
	@goose -dir db/migrations create ${name} sql

fresh: ## drop and run all containers
	$(info drop and run all containers)
	@docker-compose down
	@docker-compose up -d

sqlc:## generate slqc
	sqlc generate
test: ## run test for all
	@go test -v -cover -count=1 ./...
