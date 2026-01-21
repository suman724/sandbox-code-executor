package middleware

import (
	"net/http"
)

const sessionTokenHeader = "X-Session-Token"

func AuthBypassMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func AuthTokenMiddleware(token string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if token == "" {
				http.Error(w, "auth token not configured", http.StatusUnauthorized)
				return
			}
			if r.Header.Get(sessionTokenHeader) != token {
				http.Error(w, "invalid session token", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
