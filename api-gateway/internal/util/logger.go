package util

import (
	"context"
	"log/slog"
)

func LoggerFromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value("logger").(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}
