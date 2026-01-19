package tools

import "net/http"

type ArtifactsTool struct{}

func (ArtifactsTool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
