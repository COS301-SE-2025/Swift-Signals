PROTOC_VERSION=29.3
PROTOC_GEN_GO_VERSION=v1.36.7
PROTOC_GEN_GO_GRPC_VERSION=v1.5.1
PYTHON ?= python3
GEN_DIR=../gen
VENV_DIR=.venv

.PHONY: all install-protoc generate clean-gen generate-go generate-py clean

workflow: install-protoc generate

user-service: install-protoc-go clean-gen generate-user-service

intersection-service: install-protoc-go clean-gen generate-intersection-service

api-gateway: install-protoc-go clean-gen generate-go

generate: clean-gen generate-go generate-py

clean-gen:
	rm -rf $(GEN_DIR)
	mkdir -p $(GEN_DIR)

clean: clean-gen
	rm -rf $(VENV_DIR)

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
	. $(VENV_DIR)/bin/activate && $(PYTHON) -m grpc_tools.protoc -I .. \
		--python_out=$(GEN_DIR) \
		--pyi_out=$(GEN_DIR) \
		--grpc_python_out=$(GEN_DIR) \
		swiftsignals/common/v1/types.proto \
		swiftsignals/common/v1/simulation.proto \
		swiftsignals/simulation/v1/simulation.proto \
		swiftsignals/optimisation/v1/optimisation.proto

generate-py-dockerfile: 
	$(PYTHON) -m grpc_tools.protoc -I .. \
		--python_out=$(GEN_DIR) \
		--pyi_out=$(GEN_DIR) \
		--grpc_python_out=$(GEN_DIR) \
		swiftsignals/common/v1/types.proto \
		swiftsignals/common/v1/simulation.proto \
		swiftsignals/simulation/v1/simulation.proto \
		swiftsignals/optimisation/v1/optimisation.proto

generate-user-service:  
	protoc -I .. \
		--go_out=$(GEN_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(GEN_DIR) --go-grpc_opt=paths=source_relative \
		swiftsignals/user/v1/user.proto ; \
	cd $(GEN_DIR) && go mod init github.com/COS301-SE-2025/Swift-Signals/protos/gen && go mod tidy

generate-intersection-service:  
	protoc -I .. \
		--go_out=$(GEN_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(GEN_DIR) --go-grpc_opt=paths=source_relative \
		swiftsignals/common/v1/types.proto \
		swiftsignals/common/v1/simulation.proto \
		swiftsignals/intersection/v1/intersection.proto ; \
	cd $(GEN_DIR) && go mod init github.com/COS301-SE-2025/Swift-Signals/protos/gen && go mod tidy

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
	if [ ! -d "$(VENV_DIR)" ]; then \
		$(PYTHON) -m venv $(VENV_DIR); \
	fi
	. $(VENV_DIR)/bin/activate && pip install --upgrade pip grpcio grpcio-tools
	@echo "protoc, Go plugins, and Python gRPC tools installed!"

install-protoc-go:
	@echo "Installing protoc and gRPC tools..."
	apk add --no-cache protobuf protobuf-dev && \
    rm -rf /var/cache/apk/* && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION) && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION) 

install-protoc-py:
	@echo "Installing protoc and Python gRPC tools..."
	apk add --no-cache python3 py3-pip protobuf protobuf-dev && \
	rm -rf /var/cache/apk/* && \
	pip install --break-system-packages --upgrade pip grpcio grpcio-tools
	@echo "protoc and Python gRPC tools installed!"
