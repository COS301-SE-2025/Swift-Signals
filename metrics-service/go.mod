module github.com/COS301-SE-2025/Swift-Signals/metrics-service

replace (
	github.com/COS301-SE-2025/Swift-Signals/protos/gen => ../protos/gen
	github.com/COS301-SE-2025/Swift-Signals/shared => ../shared
)

go 1.24.3
