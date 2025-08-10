MOCKERY_VERSION=v2.53.4
GOLANGCI_LINT_VERSION=v1.62.2

.PHONY: all build clean test deps tools mocks proto lint

all: deps tools proto mocks test lint

deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

tools:
	@echo "Installing tools..."
	go get github.com/vektra/mockery/v2@$(MOCKERY_VERSION)
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

proto:
	@echo "Generating proto files..."
	cd ../protos && make install-protoc
	cd ../protos && make generate

mocks: tools
	@echo "Generating mocks with mockery $(MOCKERY_VERSION)..."
	@echo "Cleaning existing mocks..."
	rm -rf internal/mocks
	mkdir -p internal/mocks/{db,service,grpc}
	@echo "Generating mocks..."
	go run github.com/vektra/mockery/v2@$(MOCKERY_VERSION) --version
	go run github.com/vektra/mockery/v2@$(MOCKERY_VERSION)

lint: tools mocks
	@echo "Running golangci-lint $(GOLANGCI_LINT_VERSION)..."
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION) run

test: mocks
	@echo "Running tests..."
	go test -v ./internal/db/test ./internal/service/test ./internal/handler/test

test-coverage: mocks
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
