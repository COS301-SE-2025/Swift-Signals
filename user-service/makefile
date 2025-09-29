MOCKERY_VERSION=v2.53.4
GOLANGCI_LINT_VERSION=v1.62.2

.PHONY: all build clean test deps deps-post tools mocks proto lint

workflow: tools deps-pre mocks deps-post test lint

dockerfile: tools deps-pre mocks deps-post build

dockerfile-production: tools deps-pre deps-post build

deps-pre:
	@echo "Installing base dependencies..."
	go mod download

deps-post: 
	@echo "Tidying dependencies after code generation..."
	go mod tidy

tools:
	@echo "Installing tools..."
	go install github.com/vektra/mockery/v2@$(MOCKERY_VERSION)
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

mocks: 
	@echo "Generating mocks with mockery $(MOCKERY_VERSION)..."
	@echo "Cleaning existing mocks..."
	rm -rf internal/mocks
	mkdir -p internal/mocks/{db,service,grpc}
	@echo "Generating mocks..."
	mockery --version
	mockery

lint: 
	@echo "Running golangci-lint $(GOLANGCI_LINT_VERSION)..."
	golangci-lint run

test: 
	@echo "Running tests..."
	go test -v ./internal/db/test ./internal/service/test ./internal/handler/test

test-coverage: 
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out -coverpkg=./internal/db/,./internal/service/,./internal/handler/ ./internal/db/test ./internal/service/test ./internal/handler/test
	go tool cover -html=coverage.out 

clean:
	@echo "Cleaning..."
	rm -rf internal/mocks
	rm -f coverage.out coverage.html

build: 
	@echo "Building..."
	go build -o bin/user-service ./cmd/
