module github.com/COS301-SE-2025/Swift-Signals/api-gateway

go 1.24.4

replace (
	github.com/COS301-SE-2025/Swift-Signals/protos/gen => ../protos/gen
	github.com/COS301-SE-2025/Swift-Signals/shared => ../shared
)

require (
	github.com/COS301-SE-2025/Swift-Signals/protos/gen v0.1.0
	github.com/COS301-SE-2025/Swift-Signals/shared v0.0.0-00010101000000-000000000000
	github.com/swaggo/http-swagger v1.3.4
	github.com/swaggo/swag v1.16.4
	google.golang.org/grpc v1.73.0
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.20.0 // indirect
	github.com/go-openapi/spec v0.20.6 // indirect
	github.com/go-openapi/swag v0.19.15 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/swaggo/files v0.0.0-20220610200504-28940afbdbfe // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	golang.org/x/tools v0.21.1-0.20240508182429-e35e4ccd0d2d // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250324211829-b45e905df463 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
