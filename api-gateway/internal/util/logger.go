package util

import (
	"context"
	"log/slog"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/middleware"
)

func LoggerFromContext(ctx context.Context) *slog.Logger {
	if logger, ok := middleware.LoggerFromContext(ctx); ok {
		return logger
	}
	return slog.Default()
}
