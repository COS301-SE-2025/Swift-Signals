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
	cd gen/ && go mod init github.com/COS301-SE-2025/Swift-Signals/protos/gen && go mod tidy

install-protoc:
  sudo apt-get update && sudo apt-get install -y protobuf-compiler
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.0
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
  echo "${HOME}/go/bin" >> $GITHUB_PATH
