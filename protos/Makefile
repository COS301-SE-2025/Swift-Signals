PROTOC_VERSION=29.3
PROTOC_GEN_GO_VERSION=v1.36.0
PROTOC_GEN_GO_GRPC_VERSION=v1.3.0

.PHONY: install-protoc

generate: user.proto intersection.proto simulation.proto
	mkdir -p gen/user gen/intersection gen/simulation
	protoc --go_out=./gen/user/ --go_opt=paths=source_relative \
				 --go-grpc_out=./gen/user/ --go-grpc_opt=paths=source_relative \
				 user.proto
	protoc --go_out=./gen/intersection/ --go_opt=paths=source_relative \
				 --go-grpc_out=./gen/intersection/ --go-grpc_opt=paths=source_relative \
				 intersection.proto
	protoc --go_out=./gen/simulation/ --go_opt=paths=source_relative \
				 --go-grpc_out=./gen/simulation/ --go-grpc_opt=paths=source_relative \
				 simulation.proto
	protoc --go_out=./gen/optimisation/ --go_opt=paths=source_relative \
				 --go-grpc_out=./gen/optimisation/ --go-grpc_opt=paths=source_relative \
				 optimisation.proto
	cd gen/ && go mod init github.com/COS301-SE-2025/Swift-Signals/protos/gen && go mod tidy

install-protoc:
	@echo "Installing protoc $(PROTOC_VERSION)..."
	curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/protoc-$(PROTOC_VERSION)-linux-x86_64.zip
	sudo unzip -o protoc-$(PROTOC_VERSION)-linux-x86_64.zip -d /usr/local
	go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION)
	echo "${HOME}/go/bin" >> $$GITHUB_PATH
