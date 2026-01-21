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

func RequireTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if TokenFromRequest(r) == "" {
			http.Error(w, "missing session token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func TokenFromRequest(r *http.Request) string {
	return r.Header.Get(sessionTokenHeader)
}
