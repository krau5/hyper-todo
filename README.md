### hyper-todo

A simple yet complete todo API written in Go, allowing users to register accounts and manage their own tasks. The project also provides Swagger documentation for easy API interaction.

### Getting started

Clone the repo, install the dependencies, copy `.env.example` contents to `.env` and fill in the required fields
```
~ git clone git@github.com:krau5/hyper-todo.git
~ go mod download
~ cp .env.example .env
```

**The application can be launched in two ways:**

- `make dev` - run in development mode (API with Air for hot-reload, other services in Docker)
- `make prod` - run everything in Docker (production mode)

### Features
- Gin for simplified http routing
- Gorm for DB interactions
- Project architecture inspired by [go-clean-arch](https://github.com/bxcodec/go-clean-arch)
- Cookie-based JWT authentication for simplicity and security
- Swag to generate RESTful API documentation with Swagger 2.0.
- Github Actions for CI

### Scripts
- `make build` - compiles the application
- `make test` - runs all the tests
- `make run` - builds and runs the application in release mode
- `make generate-swagger` - parses annotations to generate Swagger specifications
- `make dev` - launches development environment (API with Air for hot-reload, other services in Docker)
- `make prod` - launches production environment with all services running in Docker
- `make dev-down` - stops the development environment
- `make prod-down` - stops the production environment