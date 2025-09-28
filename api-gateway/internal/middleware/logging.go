package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

type ctxKey struct{}

var loggerKey = ctxKey{}

func SetLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func LoggerFromContext(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(loggerKey).(*slog.Logger)
	if !ok {
		return slog.Default()
	}
	return logger
}

func Logging(baseLogger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = uuid.NewString()
			}
			w.Header().Set("X-Request-ID", requestID)

			requestLogger := baseLogger.With(
				"request_id", requestID,
				"method", r.Method,
				"path", r.URL.Path,
			)

			requestLogger.Info("request started")

			ctx := SetLogger(r.Context(), requestLogger)
			r = r.WithContext(ctx)

			wrapped := &wrappedWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(wrapped, r)

			requestLogger.Info("request completed",
				"status_code", wrapped.statusCode,
				"duration_ms", time.Since(start).Milliseconds(),
			)
		})
	}
}
