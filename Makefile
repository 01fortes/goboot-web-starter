.PHONY: test build examples clean

# Default target
all: test build

# Build the library
build:
	go build -v ./...

# Run tests with coverage
test:
	go test -v -cover ./...

# Run tests and generate coverage report
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Build examples
examples:
	go build -o bin/basic-example ./examples/basic
	go build -o bin/rest-api-example ./examples/rest-api

# Run the basic example
run-basic-example:
	go run ./examples/basic/main.go

# Run the REST API example
run-rest-api-example:
	go run ./examples/rest-api/main.go

# Clean build artifacts
clean:
	rm -rf bin coverage.out coverage.html

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	go vet ./...

# Update dependencies
deps-update:
	go get -u ./...
	go mod tidy

# Initialize a new project that uses this starter
init-project:
	@echo "Enter your new project name (e.g. github.com/username/project):"
	@read project_name; \
	mkdir -p $$project_name; \
	cd $$project_name; \
	go mod init $$project_name; \
	go get github.com/01fortes/goboot; \
	go get github.com/01fortes/goboot-web-starter; \
	mkdir -p cmd/server; \
	touch cmd/server/main.go; \
	echo "Project initialized at $$project_name"