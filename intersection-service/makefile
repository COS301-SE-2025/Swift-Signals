MOCKERY_VERSION=v2.53.4
GOLANGCI_LINT_VERSION=v1.62.2

.PHONY: all build clean test deps tools mocks lint

workflow: tools deps-pre mocks deps-post unit-test integration-test lint

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
	mockery --config .mockery.yml

lint: tools deps-post
	@echo "Running golangci-lint $(GOLANGCI_LINT_VERSION)..."
	golangci-lint run

unit-test: mocks deps-post
	@echo "Running unit tests..."
	go test -v ./internal/db/test ./internal/service/test ./internal/handler/test

unit-test-coverage: mocks deps-post
	@echo "Running unit tests with coverage..."
	go test -v -coverprofile=coverage.out -coverpkg=./internal/db/,./internal/service/,./internal/handler/ ./internal/db/test ./internal/service/test ./internal/handler/test
	go tool cover -html=coverage.out 

integration-test: 
	@echo "Running integration tests..."
	go test -v ./internal/test 

integration-test-coverage: 
	@echo "Running integration tests with coverage..."
	go test -v -coverprofile=coverage.out -coverpkg=./internal/db/,./internal/service/,./internal/handler/ ./internal/test
	go tool cover -html=coverage.out 

clean:
	@echo "Cleaning..."
	rm -rf internal/mocks
	rm -f coverage.out coverage.html

build:
	@echo "Building..."
	go build -o bin/intersection-service ./cmd/

