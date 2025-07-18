package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

func CreateStack(mws ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(mws) - 1; i >= 0; i-- {
			mw := mws[i]
			next = mw(next)
		}
		return next
	}
}
