PROTOC_VERSION=29.3
PROTOC_GEN_GO_VERSION=v1.36.7
PROTOC_GEN_GO_GRPC_VERSION=v1.5.1

.PHONY: install-protoc

generate: clean-gen generate-go generate-py

clean-gen:
	rm -rf ./gen/
	mkdir -p gen/user gen/intersection gen/simulation gen/optimisation

generate-go: user.proto intersection.proto simulation.proto optimisation.proto
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

generate-py: simulation.proto optimisation.proto
	python -m grpc_tools.protoc -I . \
				 --python_out=./gen/simulation/ \
				 --pyi_out=./gen/simulation/ \
				 --grpc_python_out=./gen/simulation/ \
				 simulation.proto
	touch gen/simulation/__init__.py
	python -m grpc_tools.protoc -I . \
				 --python_out=./gen/optimisation/ \
				 --pyi_out=./gen/optimisation/ \
				 --grpc_python_out=./gen/optimisation/ \
				 optimisation.proto
	touch gen/optimisation/__init__.py

install-protoc:
	@echo "Installing protoc $(PROTOC_VERSION)..."
	curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/protoc-$(PROTOC_VERSION)-linux-x86_64.zip
	sudo unzip -o protoc-$(PROTOC_VERSION)-linux-x86_64.zip -d /usr/local
	go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION)
	python -m pip install --upgrade pip
	python -m pip install grpcio grpcio-tools
	echo "${HOME}/go/bin" >> $$GITHUB_PATH
