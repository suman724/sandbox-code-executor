package tools

import "net/http"

type WorkflowsTool struct{}

func (WorkflowsTool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
