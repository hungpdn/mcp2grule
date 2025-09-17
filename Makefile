BINARY_NAME=mcp2grule
DOCKER_IMAGE=mcp2grule:latest
COVERAGE_OUT=coverage.out

# Versioning - gets the latest git tag and commit hash
VERSION ?= $(shell git describe --tags --abbrev=0)
GIT_COMMIT=$(shell git rev-parse --short HEAD)

# Linker flags to embed version information into the binary
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Commit=$(GIT_COMMIT)"

# Default target: builds the project.
.PHONY: all
all: test build

# Build the Go application
.PHONY: build
build:
	@echo "Building $(BINARY_NAME) version $(VERSION) (commit: $(COMMIT_HASH))..."
	go build -v $(LDFLAGS) -o $(BINARY_NAME) main.go

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v ./... -coverprofile=$(COVERAGE_OUT)

# Generate HTML report from coverage
.PHONY: test-coverage-html  
test-coverage-html: test-coverage
	@echo "Generating HTML coverage report..."
	go tool cover -html=$(COVERAGE_OUT)

# Clean up build artifacts
.PHONY: clean
clean:
	@echo "Cleaning up..."
	go clean
	rm -rf $(BINARY_NAME) $(COVERAGE_OUT)

# Run golangci-lint
.PHONY: lint
lint:
	@echo "Running golangci-lint..."
	which golangci-lint > /dev/null || (echo "Please install golangci-lint: https://golangci-lint.run/usage/install/"; exit 1)
	golangci-lint run

# Run MCP Inspector
.PHONY: mcp-inspector
mcp-inspector:
	@echo "Running MCP Inspector..."
	which npx > /dev/null || (echo "Please install Node.js and npm: https://nodejs.org/"; exit 1)
	npx @modelcontextprotocol/inspector

# Build Docker image
.PHONY: docker-build
docker-build:
	@echo "==> Building Docker image..."
	which docker > /dev/null || (echo "Please install Docker: https://docs.docker.com/get-docker/"; exit 1)
	docker build -t $(DOCKER_IMAGE) .

# Tidy and download Go module dependencies.
.PHONY: deps
deps:
	@echo "Tidying and downloading dependencies..."
	go mod tidy
	go mod download

# A self-documenting help target
.PHONY: help
help:
	@echo "Makefile commands:"
	@echo "  make all                - Build the project (default)"
	@echo "  make build              - Build the Go application"
	@echo "  make test               - Run tests"
	@echo "  make test-coverage      - Run tests with coverage"
	@echo "  make test-coverage-html - Generate HTML report from coverage"
	@echo "  make clean              - Clean up build artifacts"
	@echo "  make lint               - Run golangci-lint"
	@echo "  make mcp-inspector      - Run MCP Inspector"
	@echo "  make docker-build       - Build Docker image"
	@echo "  make deps               - Tidy and download Go module dependencies"
	@echo "  make help               - Show this help message"
