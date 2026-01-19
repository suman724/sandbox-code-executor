package tools

import "net/http"

type JobsTool struct{}

func (JobsTool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
