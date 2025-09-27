PROTOC_VERSION=29.3
PROTOC_GEN_GO_VERSION=v1.36.7
PROTOC_GEN_GO_GRPC_VERSION=v1.5.1
PYTHON ?= python3
GEN_DIR=../gen

.PHONY: install-protoc generate clean-gen generate-go generate-py

all: install-protoc generate

generate: clean-gen generate-go generate-py

clean-gen:
	rm -rf $(GEN_DIR)
	mkdir -p $(GEN_DIR)

generate-go:
	protoc -I .. \
		--go_out=$(GEN_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(GEN_DIR) --go-grpc_opt=paths=source_relative \
		swiftsignals/common/v1/types.proto \
		swiftsignals/common/v1/simulation.proto \
		swiftsignals/user/v1/user.proto \
		swiftsignals/intersection/v1/intersection.proto \
		swiftsignals/simulation/v1/simulation.proto \
		swiftsignals/optimisation/v1/optimisation.proto
	cd $(GEN_DIR) && go mod init github.com/COS301-SE-2025/Swift-Signals/protos/gen && go mod tidy

generate-py:
	$(PYTHON) -m grpc_tools.protoc -I .. \
		--python_out=$(GEN_DIR) \
		--pyi_out=$(GEN_DIR) \
		--grpc_python_out=$(GEN_DIR) \
		swiftsignals/common/v1/types.proto \
		swiftsignals/simulation/v1/simulation.proto \
		swiftsignals/optimisation/v1/optimisation.proto

install-protoc:
	@echo "Installing protoc and gRPC tools..."
ifeq ($(shell uname),Darwin)
	brew list protobuf >/dev/null 2>&1 || brew install protobuf
else
	curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/protoc-$(PROTOC_VERSION)-linux-x86_64.zip
	sudo unzip -o protoc-$(PROTOC_VERSION)-linux-x86_64.zip -d /usr/local
endif
	go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION)
	export PATH=$$PATH:$(shell go env GOPATH)/bin
	if [ ! -d ".venv" ]; then \
		python3 -m venv .venv; \
	fi
	. .venv/bin/activate && pip install --upgrade pip grpcio grpcio-tools
	@echo "protoc, Go plugins, and Python gRPC tools installed!"
