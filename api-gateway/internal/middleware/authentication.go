package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/util"
	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"github.com/COS301-SE-2025/Swift-Signals/shared/jwt"
)

type contextKey string

const userIDKey contextKey = "userID"

func AuthMiddleware(secret string, paths ...string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			excluded := NewPathSet(paths...)
			if excluded.Contains(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			logger := LoggerFromContext(r.Context())
			logger.Info("authentication user request")

			token, err := util.GetToken(r)
			if err != nil {
				logger.Error("failed to retrieve authorization token",
					"error", err.Error())
				util.SendErrorResponse(
					w,
					errs.NewUnauthorizedError(
						fmt.Sprintf("failed to authorize: %s", err.Error()),
						map[string]any{"error": err},
					),
				)
				return
			}

			jwt.Init([]byte(secret))
			claims, err := jwt.ParseToken(token)
			if err != nil {
				logger.Error("failed to parse token",
					"error", err.Error())
				util.SendErrorResponse(
					w,
					errs.NewUnauthorizedError("invalid token", map[string]any{}),
				)
				return
			}

			userID := claims.UserID
			if userID == "" {
				logger.Warn("user ID missing in jwt")
				util.SendErrorResponse(
					w,
					errs.NewUnauthorizedError("user ID missing in jwt", map[string]any{}),
				)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

type PathSet map[string]struct{}

func NewPathSet(paths ...string) PathSet {
	ps := make(PathSet)
	for _, p := range paths {
		ps[p] = struct{}{}
	}
	return ps
}

func (ps PathSet) Contains(path string) bool {
	_, exists := ps[path]
	return exists
}

func GetUserID(r *http.Request) (string, bool) {
	userID, ok := r.Context().Value(userIDKey).(string)
	return userID, ok
}
