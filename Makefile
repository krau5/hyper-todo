build:
	@go build -o bin/hyper-todo cmd/api/main.go

test:
	@go test -v ./...

run: build
	@./bin/hyper-todo