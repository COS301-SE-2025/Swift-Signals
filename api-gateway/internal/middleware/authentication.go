package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/COS301-SE-2025/Swift-Signals/api-gateway/internal/util"
	errs "github.com/COS301-SE-2025/Swift-Signals/shared/error"
	"github.com/COS301-SE-2025/Swift-Signals/shared/jwt"
)

type contextKey string

const (
	userIDKey contextKey = "userID"
	roleKey   contextKey = "role"
)

func AuthMiddleware(secret string, paths ...string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			excluded := NewPathSet(paths...)
			if excluded.Contains(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			logger := LoggerFromContext(r.Context())
			logger.Info("authenticating user request...")

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

			ctx := SetUserID(r.Context(), userID)
			ctx = SetRole(ctx, claims.Role)
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
	for p := range ps {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}

func SetUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func SetRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, roleKey, role)
}

func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey).(string)
	return userID, ok
}

func GetRole(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(roleKey).(string)
	return role, ok
}
