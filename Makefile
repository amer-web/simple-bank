name ?= ''
.PHONY: help migrate-up migrate-down fresh migrate-create sqlc test run mock proto
help: ## Print help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
migrate-up: ## run migrate up
	$(info migrate up running....)
#	@migrate -path db/migration -database "postgresql://amer:amer@127.0.0.1:5432/simple_bank?sslmode=disable" force 1
#	@migrate -path db/migration -database "postgresql://amer:amer@127.0.0.1:5432/simple_bank?sslmode=disable" -verbose up
	@goose -dir db/migrations postgres "postgresql://amer:amer@127.0.0.1:5432/simple_bank?sslmode=disable" up
migrate-down: ## run migrate down by one
	$(info migrate down running....)
	@goose -dir db/migrations postgres "postgresql://amer:amer@127.0.0.1:5432/simple_bank?sslmode=disable" down
migrate-reset: ## run migrate reset
	$(info migrate reset running....)
	@goose -dir db/migrations postgres "postgresql://amer:amer@127.0.0.1:5432/simple_bank?sslmode=disable" reset


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

run: ## run server
	@go run main.go

mock: ## run mock
	@mockgen -package mockdb -destination db/mock/store.go github.com/amer-web/simple-bank/db/sqlc Store
proto: ## generate go code from proto files
	rm -r pb/*
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
        --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
        --grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
        --validate_out="lang=go:pb" --validate_opt=paths=source_relative \
        --experimental_allow_proto3_optional \
        proto/*.proto
evans: ## run evans
	evans -r repl -p 50051
