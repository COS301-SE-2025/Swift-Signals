//go:build tools
// +build tools

package tools

import (
	_ "github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/mocks/client"
	_ "github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/mocks/grpc_client"
	_ "github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/mocks/service"
	_ "github.com/vektra/mockery/v2"
)
