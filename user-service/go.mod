module github.com/COS301-SE-2025/Swift-Signals/user-service

go 1.24.4

replace (
	github.com/COS301-SE-2025/Swift-Signals/protos/gen => ../protos/gen
	github.com/COS301-SE-2025/Swift-Signals/shared => ../shared
)

toolchain go1.24.5

require (
	github.com/lib/pq v1.10.9
	github.com/stretchr/testify v1.10.0
	golang.org/x/crypto v0.36.0
	google.golang.org/grpc v1.73.0
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250324211829-b45e905df463 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/COS301-SE-2025/Swift-Signals/protos/gen v0.0.0-00010101000000-000000000000
	github.com/COS301-SE-2025/Swift-Signals/shared v0.0.0-00010101000000-000000000000
	github.com/DATA-DOG/go-sqlmock v1.5.2
	github.com/google/uuid v1.6.0
)
