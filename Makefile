build:
	@go build -o bin/hyper-todo cmd/api/main.go

test:
	@go test -v ./...

run: build
	@export GIN_MODE=release && ./bin/hyper-todo

generate-swagger:
	@swag init --parseDependency --parseInternal -g cmd/api/main.go