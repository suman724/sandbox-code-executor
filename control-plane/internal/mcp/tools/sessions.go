package tools

import "net/http"

type SessionsTool struct{}

func (SessionsTool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
