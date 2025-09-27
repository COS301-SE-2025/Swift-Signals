MOCKERY_VERSION=v2.53.4
GOLANGCI_LINT_VERSION=v1.62.2

.PHONY: all build clean test deps tools mocks lint

all: tools mocks deps unit-test lint

deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

tools:
	@echo "Installing tools..."
	go install github.com/vektra/mockery/v2@$(MOCKERY_VERSION)
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

mocks: tools
	@echo "Generating mocks with mockery $(MOCKERY_VERSION)..."
	@echo "Cleaning existing mocks..."
	rm -rf internal/mocks
	mkdir -p internal/mocks/{client,service,grpc_client}
	@echo "Generating mocks..."
	mockery --version
	mockery --config .mockery.yml

lint: tools
	@echo "Running golangci-lint $(GOLANGCI_LINT_VERSION)..."
	golangci-lint run

unit-test: mocks
	@echo "Running unit tests..."
	go test -v ./internal/client/test/intersection ./internal/client/test/optimisation ./internal/client/test/user
	go test -v ./internal/service/test/admin ./internal/service/test/auth ./internal/service/test/intersection ./internal/service/test/profile
	go test -v ./internal/handler/test/admin ./internal/handler/test/auth ./internal/handler/test/intersection ./internal/handler/test/profile

clean:
	@echo "Cleaning..."
	rm -rf internal/mocks
	rm -f coverage.out coverage.html

build:
	@echo "Building..."
	go build -o bin/api-gateway ./cmd/


