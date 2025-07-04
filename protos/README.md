# Shared Protofile Contracts for Microservices
This folder is dedicated for specifying the protocols that the gRPC servers in our all of the Swift-Signals microservices will follow.

## Generating Code
To generate code, you will first need to install the `protoc` tool.
Run all of the following commands from inside the protos directory.


### To update simulation.proto generated code

```bash
python -m grpc_tools.protoc -I . \
       --python_out=./gen/simulation/ \
       --pyi_out=./gen/simulation/ \
       --grpc_python_out=./gen/simulation/ \
       simulation.proto
```


### To update user.proto generated code
```bash
protoc --go_out=./gen/user/ --go_opt=paths=source_relative \
       --go-grpc_out=./gen/user/ --go-grpc_opt=paths=source_relative \
       user.proto
```


### To update intersection.proto generated code
```bash
protoc --go_out=./gen/intersection/ --go_opt=paths=source_relative \
       --go-grpc_out=./gen/intersection/ --go-grpc_opt=paths=source_relative \
       intersection.proto
```

### To update simulation.proto generated code
```bash
protoc --go_out=./gen/simulation/ --go_opt=paths=source_relative \
       --go-grpc_out=./gen/simulation/ --go-grpc_opt=paths=source_relative \
       simulation.proto
```
