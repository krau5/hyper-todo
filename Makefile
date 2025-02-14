.PHONY: build test run generate-swagger dev prod dev-down prod-down

build:
	@go build -o bin/hyper-todo cmd/api/main.go

test:
	@go test -v ./...

run: build
	@export GIN_MODE=release && ./bin/hyper-todo

generate-swagger:
	@swag init --parseDependency --parseInternal -g cmd/api/main.go

dev:
	@docker compose -f docker-compose.dev.yml up -d
	@air

prod:
	@docker compose -f docker-compose.prod.yml up -d

dev-down:
	@docker compose -f docker-compose.dev.yml down

prod-down:
	@docker compose -f docker-compose.prod.yml down
